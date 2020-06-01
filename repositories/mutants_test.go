package repositories_test

import (
	mgo "github.com/sebastianreh/mutants-detector/dal/mongo"
	rds "github.com/sebastianreh/mutants-detector/dal/redis"
	"github.com/sebastianreh/mutants-detector/models"
	"github.com/sebastianreh/mutants-detector/repositories"
	mongoMock "github.com/sebastianreh/mutants-detector/test/mocks/dal/mongo"
	redisMock "github.com/sebastianreh/mutants-detector/test/mocks/dal/redis"
	cacheMock "github.com/sebastianreh/mutants-detector/test/mocks/repositories"
	"context"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

var _ = Describe("Mutants", func() {

	var ctrl *gomock.Controller
	var mockMongoDatabase *mongoMock.MockIMongoDatabase
	var mockRedisDatabase *redisMock.MockIRedisDatabase
	var mockCache *cacheMock.MockIMutantCache
	var dnaId string
	var mutantsCache *models.MutantsStats
	var subject models.Subject
	var mutantReposistory repositories.IMutantRepository
	var redisClient *redis.Client
	var mongoDatabase *mongo.Database

	BeforeSuite(func() {
		redisClient = rds.NewRedisTestDatabase()
		mongoDatabase = mgo.BuildTestMongoDatabase()
		mgo.ConfigureIndexes(mongoDatabase)
	})

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockMongoDatabase = mongoMock.NewMockIMongoDatabase(ctrl)
		mockRedisDatabase = redisMock.NewMockIRedisDatabase(ctrl)
		mockCache = cacheMock.NewMockIMutantCache(ctrl)
		mutantReposistory = repositories.NewMutantRepository(mockMongoDatabase, mockRedisDatabase, mockCache)
		redisClient.SAdd(repositories.HumansSavedKey, "4CTTACGGGCCCCCGCT")
		redisClient.SAdd(repositories.HumansNotSavedKey, "4TGATTTTTGTAATATT")
		redisClient.SAdd(repositories.MutantsSavedKey, "4GGGACCCGGAACGGTA")
		redisClient.SAdd(repositories.MutantsNotSavedKey, "4CTACATCTCGAGCTGA")
	})

	AfterEach(func() {
		redisClient.SRem(repositories.HumansSavedKey, "4CTTACGGGCCCCCGCT")
		redisClient.SRem(repositories.HumansNotSavedKey, "4TGATTTTTGTAATATT")
		redisClient.SRem(repositories.MutantsSavedKey, "4GGGACCCGGAACGGTA")
		redisClient.SRem(repositories.MutantsNotSavedKey, "4CTACATCTCGAGCTGA")
		ctrl.Finish()
	})

	Context("when a mutant or human request is not in database", func() {
		It("should save mutant in database", func() {
			dnaId = "4CTTACGGGCCCCCGCT"
			subject = models.Subject{
				Id:       dnaId,
				Dna:      []string{"CTTA", "CGGG", "CCCC", "CGCT"},
				IsMutant: true,
			}
			newMutantsCache := models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt64(1, 100)),
				CountHumanDna:  float64(faker.RandomInt64(1, 100)),
			}
			newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
			mutantsCache = &newMutantsCache
			mockRedisDatabase.EXPECT().RedisClient().Return(redisClient).AnyTimes()
			mockMongoDatabase.EXPECT().SubjectsDatabase().Return(mongoDatabase).AnyTimes()
			mockCache.EXPECT().GetStatsFromCache().Return(mutantsCache)
			mockCache.EXPECT().SaveStatsInCache(mutantsCache).AnyTimes()
			mutantReposistory.SaveSubjectIteration(subject)
		})

		It("should save human in database", func() {
			dnaId = "4CTATTACGACCCTCGA"
			subject = models.Subject{
				Id:       dnaId,
				Dna:      []string{"CTAT", "TACG", "ACCC", "TCGA"},
				IsMutant: false,
			}
			newMutantsCache := models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt64(1, 100)),
				CountHumanDna:  float64(faker.RandomInt64(1, 100)),
			}
			newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
			mutantsCache = &newMutantsCache
			mockRedisDatabase.EXPECT().RedisClient().Return(redisClient).AnyTimes()
			mockMongoDatabase.EXPECT().SubjectsDatabase().Return(mongoDatabase).AnyTimes()
			mockCache.EXPECT().GetStatsFromCache().Return(mutantsCache)
			mockCache.EXPECT().SaveStatsInCache(mutantsCache).AnyTimes()
			mutantReposistory.SaveSubjectIteration(subject)
		})

		Context("when a mutant or human request is already in database", func() {
			It("should not save mutant in database", func() {
				redisClient.SAdd(repositories.MutantsSavedKey, "4CTTACGGGCCCCCGCT")
				dnaId = "4CTTACGGGCCCCCGCT"
				subject = models.Subject{
					Id:       dnaId,
					Dna:      []string{"CTTA", "CGGG", "CCCC", "CGCT"},
					IsMutant: true,
				}
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(faker.RandomInt64(1, 100)),
					CountHumanDna:  float64(faker.RandomInt64(1, 100)),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockRedisDatabase.EXPECT().RedisClient().Return(redisClient).AnyTimes()
				mockMongoDatabase.EXPECT().SubjectsDatabase().Return(mongoDatabase).AnyTimes()
				mockCache.EXPECT().GetStatsFromCache().Return(mutantsCache)
				mockCache.EXPECT().SaveStatsInCache(mutantsCache).AnyTimes()
				mutantReposistory.SaveSubjectIteration(subject)
			})

			It("should not save human in database", func() {
				redisClient.SAdd(repositories.HumansSavedKey, "4CTGCCATGGACCTGTA")
				dnaId = "4CTATTACGACCCTCGA"
				subject = models.Subject{
					Id:       dnaId,
					Dna:      []string{"CTAT", "TACG", "ACCC", "TCGA"},
					IsMutant: false,
				}
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(faker.RandomInt64(1, 100)),
					CountHumanDna:  float64(faker.RandomInt64(1, 100)),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockRedisDatabase.EXPECT().RedisClient().Return(redisClient).AnyTimes()
				mockMongoDatabase.EXPECT().SubjectsDatabase().Return(mongoDatabase).AnyTimes()
				mockCache.EXPECT().GetStatsFromCache().Return(mutantsCache)
				mockCache.EXPECT().SaveStatsInCache(mutantsCache).AnyTimes()
				mutantReposistory.SaveSubjectIteration(subject)
			})
		})

		Context("when deploying mircroservice", func() {
			It("should upload all mongoData into redis", func() {
				mutantsCollection := mongoDatabase.Collection(mgo.MutantsCollection)
				humansCollection := mongoDatabase.Collection(mgo.HumansCollection)
				mutantsSubjects := []models.Subject{
					{
						Id:       "4CTTACGGGCCCCCGCT",
						Dna:      []string{"CTTA", "CGGG", "CCCC", "CGCT"},
						IsMutant: true,
					},
					{
						Id:       "4GCTGTTCCCCTGCATT",
						Dna:      []string{"GCTG", "TTCC", "CCTG", "CATT"},
						IsMutant: true,
					},
				}
				humansSubjects := []models.Subject{
					{
						Id:       "4CTATTACGACCCTCGA",
						Dna:      []string{"CTAT", "TACG", "ACCC", "TCGA"},
						IsMutant: false,
					},
					{
						Id:       "4TCTGAGAGTAGCCCGC",
						Dna:      []string{"TCTG", "AGAG", "TAGC", "CCGC"},
						IsMutant: false,
					},
				}
				mutantsOperation := make([]mongo.WriteModel, 0)
				for _, mutant := range mutantsSubjects {
					mutantModel := repositories.BuildSubjectModelDB(mutant)
					mutantsOperation = append(mutantsOperation, mutantModel)
				}
				mutantsCollection.BulkWrite(context.Background(), mutantsOperation)
				humansOperation := make([]mongo.WriteModel, 0)
				for _, human := range humansSubjects {
					humanModel := repositories.BuildSubjectModelDB(human)
					humansOperation = append(humansOperation, humanModel)
				}
				humansCollection.BulkWrite(context.Background(), mutantsOperation)
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(3),
					CountHumanDna:  float64(4),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockRedisDatabase.EXPECT().RedisClient().Return(redisClient).AnyTimes()
				mockMongoDatabase.EXPECT().SubjectsDatabase().Return(mongoDatabase).AnyTimes()
				mockCache.EXPECT().SaveStatsInCache(gomock.Any()).AnyTimes()
				mutantReposistory.BuildRedisData()
			})
		})

		Context("when saving in cache", func() {
			It("should save into cache", func() {
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(3),
					CountHumanDna:  float64(4),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockCache.EXPECT().SaveStatsInCache(gomock.Any()).AnyTimes()
				mutantReposistory.SaveStatsInCache(mutantsCache)
			})
		})

		Context("when getting data from cache", func() {
			It("should return cache", func() {
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(3),
					CountHumanDna:  float64(4),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockCache.EXPECT().GetStatsFromCache().AnyTimes()
				mutantReposistory.GetStatsFromCache()
			})
		})

		Context("when asked for subject stats", func() {
			It("should return stats", func() {
				newMutantsCache := models.MutantsStats{
					CountMutantDna: float64(3),
					CountHumanDna:  float64(4),
				}
				newMutantsCache.Ratio = newMutantsCache.CountHumanDna / newMutantsCache.CountHumanDna
				mutantsCache = &newMutantsCache
				mockCache.EXPECT().GetStatsFromCache().AnyTimes()
				_, err := mutantReposistory.GetSubjectsStats()
				Expect(err).Should(BeNil())
			})
		})
	})
})
