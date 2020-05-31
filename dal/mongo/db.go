package mongo

import (
	. "ExamenMeLiMutante/settings"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	MongoDatabase struct {
		client *mongo.Client
	}

	IMongoDatabase interface {
		SubjectsDatabase() *mongo.Database
		ConfigureDatabase()
	}
)

func NewMongoClient() IMongoDatabase {
	return MongoDatabase{}
}

func (mgo MongoDatabase) subjectsDatabaseClient() *mongo.Client {
	if mgo.client != nil {
		return mgo.client
	}
	mgo.client = BuildMongoDatabaseClient(ProjectSettings.Database.ConnectionString)
	return mgo.client
}

func (mgo MongoDatabase) SubjectsDatabase() *mongo.Database {
	return mgo.subjectsDatabaseClient().Database(ProjectSettings.Database.MutantsDbName)
}
