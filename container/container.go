package container

import (
	"github.com/sebastianreh/mutants-detector/controllers"
	"github.com/sebastianreh/mutants-detector/dal/mongo"
	"github.com/sebastianreh/mutants-detector/dal/redis"
	"github.com/sebastianreh/mutants-detector/repositories"
	"github.com/sebastianreh/mutants-detector/services"
	"github.com/sebastianreh/mutants-detector/services/finder"
)

var (
	// controllers
	MutantController controllers.IMutantController
	// services
	MutantService       services.IMutantService
	MutantFinderService finder.IMutantFinderService
	// repositories
	MutantRepository repositories.IMutantRepository
	MutantCache      repositories.IMutantCache
	// databases
	MongoDatabase mongo.IMongoDatabase
	RedisDatabase redis.IRedisDatabase
)

// Aquí se inician e inyectan las dependencias

func init() {
	MutantFinderService = finder.NewMutantFinderService()
	MongoDatabase = mongo.NewMongoClient()
	RedisDatabase = redis.NewRedisDatabase()
	MutantCache = repositories.NewMutantCacheRepository()
	MutantRepository = repositories.NewMutantRepository(MongoDatabase, RedisDatabase, MutantCache)
	MutantService = services.NewMutantService(MutantFinderService, MutantRepository)
	MutantController = controllers.NewMutantController(MutantService)
}
