package container

import (
	"ExamenMeLiMutante/controllers"
	"ExamenMeLiMutante/dal"
	"ExamenMeLiMutante/dal/mongo"
	"ExamenMeLiMutante/repositories"
	"ExamenMeLiMutante/services"
	"ExamenMeLiMutante/services/finder"
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
	RedisDatabase dal.IRedisDatabase
)

// Aquí se inician e inyectan las dependencias

func init() {
	MutantFinderService = finder.NewMutantFinderService()
	MongoDatabase = mongo.NewMongoClient()
	RedisDatabase = dal.NewRedisDatabase()
	MutantCache = repositories.NewMutantCacheRepository()
	MutantRepository = repositories.NewMutantRepository(MongoDatabase, RedisDatabase, MutantCache)
	MutantService = services.NewMutantService(MutantFinderService, MutantRepository)
	MutantController = controllers.NewMutantController(MutantService)
}
