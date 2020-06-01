package repositories

import (
	mgo "ExamenMeLiMutante/dal/mongo"
	rds "ExamenMeLiMutante/dal/redis"
	"ExamenMeLiMutante/models"
	"ExamenMeLiMutante/utils"
	"context"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	MutantRepository struct {
		mongoDb mgo.IMongoDatabase
		redis   rds.IRedisDatabase
		cache   IMutantCache
	}
)

const (
	MutantsSavedKey    = "mutants:saved:dna"
	HumansSavedKey     = "humans:saved:dna"
	MutantsNotSavedKey = "mutants:notSaved:dna"
	HumansNotSavedKey  = "humans:notSaved:dna"
	MutantStatus       = "mutant"
	HumanStatus        = "human"
	NotFoundStatus     = "notFound"
	notSavedLimit      = 1000
	notSavedTimeLimit  = 5
)

var (
	notSavedCount int64
	upsert        = true
	saveTimeLimit = time.Time{}
)

type IMutantRepository interface {
	SaveSubjectIteration(models.Subject)
	GetSubjectStatus(dnaId string) string
	QueueDatabaseOperations(subsChan <-chan *models.Subject)
	GetSubjectsStats() (models.MutantsPreStats, error)
	GetStatsFromCache() *models.MutantsStats
	SaveStatsInCache(stats *models.MutantsStats)
	BuildRedisData()
}

var mongoDatabase *mongo.Database
var redisClient *redis.Client
var subjectCollection *mongo.Collection

func (repository MutantRepository) startRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = repository.redis.RedisClient()
	return redisClient
}

func (repository MutantRepository) startMongoClient() *mongo.Database {
	if mongoDatabase != nil {
		return mongoDatabase
	}
	mongoDatabase = repository.mongoDb.SubjectsDatabase()
	return mongoDatabase
}

func NewMutantRepository(mongoDb mgo.IMongoDatabase, redis rds.IRedisDatabase, cache IMutantCache) IMutantRepository {
	return MutantRepository{
		mongoDb: mongoDb,
		redis:   redis,
		cache:   cache,
	}
}

// Aqui se crea la cola de operaciones sobre la base de datos, para evitar abrir conexiones innecesarias

func (repository MutantRepository) QueueDatabaseOperations(subsChan <-chan *models.Subject) {
	repository.startRedisClient()
	for subject := range subsChan {
		repository.SaveSubjectIteration(*subject)
	}
}

// Aqui se guardan las iteraciones de los sujetos

func (repository MutantRepository) SaveSubjectIteration(subject models.Subject) {
	client := repository.startRedisClient()
	key := selectConditionKey(subject)
	statsCache := repository.cache.GetStatsFromCache()

	result, err := client.SAdd(key, subject.Id).Result()
	if err != nil {
		log.Errorf("repositories.SaveSubjectIteration | Error adding subject to Redis :%v", err)
	}
	if subject.IsMutant {
		statsCache.CountMutantDna = statsCache.CountMutantDna + float64(result)
	} else {
		statsCache.CountHumanDna = statsCache.CountHumanDna + float64(result)
	}
	repository.cache.SaveStatsInCache(statsCache)
	repository.handleIterationsDb(result)
}

// Es el handler de iteraciónes. Al cruzar los threasholds, revisa la data guardada en los NotSaved, transfiere la data
// a saved, genera los subjects y los guarda en mongo

func (repository MutantRepository) handleIterationsDb(result int64) {
	notSavedCount = notSavedCount + result
	if time.Now().After(saveTimeLimit) || notSavedCount > notSavedLimit {
		notSavedCount = 0
		saveTimeLimit = time.Now().Add(time.Second * notSavedTimeLimit)

		mutantsToSave := repository.transferAndGenerateSubjects(MutantStatus)
		humansToSave := repository.transferAndGenerateSubjects(HumanStatus)

		if len(mutantsToSave) != 0 {
			err := repository.saveBulkSubjectsInDb(mutantsToSave, MutantStatus)
			if err != nil {
				log.Errorf("repositories.SaveSubjectIteration | Error :%v", err)
			}
		}
		if len(humansToSave) != 0 {
			err := repository.saveBulkSubjectsInDb(humansToSave, HumanStatus)
			if err != nil {
				log.Errorf("repositories.SaveSubjectIteration | Error :%v", err)
			}
		}
	}
	notSavedCount++
}

// Aquí se guardan los sujetos en Mongo

func (repository MutantRepository) saveBulkSubjectsInDb(subjects []models.Subject, status string) error {
	switch status {
	case MutantStatus:
		subjectCollection = repository.startMongoClient().Collection(mgo.MutantsCollection)
	case HumanStatus:
		subjectCollection = repository.startMongoClient().Collection(mgo.HumansCollection)
	}
	subjectsOperations := make([]mongo.WriteModel, 0)

	for _, subject := range subjects {
		mutantModel := BuildSubjectModelDB(subject)
		subjectsOperations = append(subjectsOperations, mutantModel)
	}

	bwSubjects, err := subjectCollection.BulkWrite(context.Background(), subjectsOperations)

	if err != nil {
		log.Errorf("repositories.SaveSubjectIteration | Error saving models.Subject in database: %v", err)
	}

	if bwSubjects.UpsertedCount != 0 {
		log.Printf("repositories.SaveSubjectIteration | saved %d %ss in database", bwSubjects.UpsertedCount, status)
	}

	return err
}

// Aqui se pregunta por el status del sujeto, comparándolo contra la data existente en Redis

func (repository MutantRepository) GetSubjectStatus(dnaId string) string {
	client := repository.startRedisClient()
	mutantSavedInDb, _ := client.SIsMember(MutantsSavedKey, dnaId).Result()
	mutantNotSavedInDb, _ := client.SIsMember(MutantsNotSavedKey, dnaId).Result()
	if mutantSavedInDb || mutantNotSavedInDb {
		return MutantStatus
	}
	humanSavedInDb, _ := client.SIsMember(HumansSavedKey, dnaId).Result()
	humanNotSavedInDb, _ := client.SIsMember(HumansNotSavedKey, dnaId).Result()
	if humanSavedInDb || humanNotSavedInDb {
		return HumanStatus
	}
	return NotFoundStatus
}

// Aqui se pregunta por los stats de tosdos los sujetos, la existente en Redis

func (repository MutantRepository) GetSubjectsStats() (models.MutantsPreStats, error) {
	client := repository.startRedisClient()
	mutantsNotSaved, err := client.SMembers(MutantsNotSavedKey).Result()
	mutantsNotSavedCount := len(mutantsNotSaved)
	if err != nil {
		log.Errorf("repositories.BuildRedisData | Error reading mutants counter: %v", err)
	}
	humansNotSaved, err := client.SMembers(HumansNotSavedKey).Result()
	humansNotSavedCount := len(humansNotSaved)
	if err != nil {
		log.Errorf("repositories.BuildRedisData | Error reading humans counter: %v", err)
	}
	mutantsSaved, err := client.SMembers(MutantsSavedKey).Result()
	mutantsSavedCount := len(mutantsSaved)
	if err != nil {
		log.Errorf("repositories.BuildRedisData | Error reading mutants counter: %v", err)
	}
	humansSaved, err := client.SMembers(HumansSavedKey).Result()
	humansSavedCount := len(humansSaved)
	if err != nil {
		log.Errorf("repositories.BuildRedisData | Error reading humans counter: %v", err)
	}
	preStats := models.MutantsPreStats{
		CountMutants: mutantsNotSavedCount + mutantsSavedCount,
		CountHumans:  humansNotSavedCount + humansSavedCount,
	}
	return preStats, nil
}

// Aqui subr la data que existe en Mongo a Redis, para poder compararla rápidamente
// frente a los nuevos requests y evitar la repetición de la próxima data que se tenga que guardar

func (repository MutantRepository) BuildRedisData() {
	humansCount := repository.BuildRedisSubjectData(HumanStatus, mgo.HumansCollection)
	mutantsCount := repository.BuildRedisSubjectData(MutantStatus, mgo.MutantsCollection)
	preStats := models.MutantsPreStats{
		CountMutants: humansCount,
		CountHumans:  mutantsCount,
	}
	mutantsStats := utils.CalculateMutantStats(preStats)
	repository.cache.SaveStatsInCache(mutantsStats)
}

func (repository MutantRepository) BuildRedisSubjectData(status string, collection string) int {
	savedKey, _ := selectKeysByStatus(status)
	client := repository.startRedisClient()
	subjectsTosSave := repository.transferAndGenerateSubjects(status)
	if len(subjectsTosSave) > 0 {
		_ = repository.saveBulkSubjectsInDb(subjectsTosSave, MutantStatus)
	}
	subjectsToRemove, _ := client.SMembers(savedKey).Result()
	client.SRem(savedKey, subjectsToRemove)
	subjectsCount := repository.buildRedisDataByStatus(savedKey, collection)
	return subjectsCount
}

// Aqui sube la data que existe en Mongo a Redis, para poder compararla rápidamente
// frente a los nuevos requests y a la próxima data que debera guardar en Mongo

func (repository MutantRepository) buildRedisDataByStatus(key string, collection string) int {
	subjectCollection = repository.startMongoClient().Collection(collection)
	client := repository.redis.RedisClient()
	query := bson.M{}
	var counter int
	findOptions := options.FindOptions{
		Projection: bson.M{
			"_id": 0,},
	}
	cursor, err := subjectCollection.Find(context.TODO(), query, &findOptions)
	if err != nil {
		log.Errorf("repositories.buildRedisDataByStatus | Error finding query: %v", err)
	}

	defer func() {
		err := cursor.Close(context.Background())
		if err != nil {
			log.Errorf("repositories.buildRedisDataByStatus | Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(context.Background()) {
		subject := new(models.Subject)
		if err := cursor.Decode(subject); err != nil {
			log.Errorf("repositories.buildRedisDataByStatus | Error decoding Objects: %v", err)
		}
		_, err := client.SAdd(key, subject.Id).Result()
		if err != nil {
			log.Errorf("repositories.buildRedisDataByStatus | Error adding subject to Redis  :%v", err)
		}
		counter++
	}
	return counter
}

// Aqui transfiere los sujetos del NotSavedKey al SavedKey y arma los sujetos para el guardado en Mongo

func (repository MutantRepository) transferAndGenerateSubjects(status string) []models.Subject {
	savedKey, notSavedKey := selectKeysByStatus(status)
	client := repository.startRedisClient()
	subjects := new([]models.Subject)

	subjectsToSaveDiff, err := client.SDiff(notSavedKey, savedKey).Result()
	if err != nil {
		log.Errorf("repositories.transferAndGenerateSubjects | Error :%v", err)
	}
	if len(subjectsToSaveDiff) > 0 {
		_, err = client.SRem(notSavedKey, subjectsToSaveDiff).Result()
		if err != nil {
			log.Errorf("repositories.transferAndGenerateSubjects | Error :%v", err)
		}
	}
	subjectsNowInSaved, _ := client.SAdd(savedKey, subjectsToSaveDiff).Result()

	log.Printf("Moving %d %ss from not saved key to saved key", subjectsNowInSaved, status)

	futureSubjects := generateSubjectsAsync(generateSubjects, subjectsToSaveDiff, status)
	subjects = <-futureSubjects
	return *subjects
}

// Devuelve la cache de stats

func (repository MutantRepository) GetStatsFromCache() *models.MutantsStats {
	cache := repository.cache.GetStatsFromCache()
	return cache
}

// Guarda la cache de stats

func (repository MutantRepository) SaveStatsInCache(stats *models.MutantsStats) {
	repository.cache.SaveStatsInCache(stats)
}
