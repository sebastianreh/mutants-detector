package mongo

import (
	. "ExamenMeLiMutante/settings"
	"context"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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

var ctxTimeout = time.Duration(30) * time.Second
var dbTimeout = time.Duration(ProjectSettings.Database.Timeout) * time.Second
var dbSocketTimeout = time.Duration(60) * time.Second
var dbMaxConnIdleTime = time.Duration(ProjectSettings.Database.MaxConnIdleTime) * time.Second

func NewMongoClient() IMongoDatabase {
	return MongoDatabase{}
}

func (mgo MongoDatabase) subjectsDatabaseClient() *mongo.Client {
	if mgo.client != nil {
		return mgo.client
	}
	mgo.client = BuildMongoDatabaseClient()
	return mgo.client
}

func (mgo MongoDatabase) SubjectsDatabase() *mongo.Database {
	return mgo.subjectsDatabaseClient().Database(ProjectSettings.Database.MutantsDbName)
}

func BuildMongoDatabaseClient() *mongo.Client {
	clientOptions := new(options.ClientOptions)
	clientOptions.ApplyURI(ProjectSettings.Database.ConnectionString)
	clientOptions.SetAppName(ProjectSettings.ProjectName)
	clientOptions.SetConnectTimeout(dbTimeout)
	clientOptions.SetSocketTimeout(dbSocketTimeout)
	clientOptions.SetServerSelectionTimeout(dbTimeout)
	clientOptions.SetMaxPoolSize(ProjectSettings.Database.MaxConnections)
	clientOptions.SetMaxConnIdleTime(dbMaxConnIdleTime)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Errorf("BuildMongoDatabaseClient | Error connecting to DB: %s", err.Error())
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), ctxTimeout)

	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("BuildMongoDatabaseClient | Error connecting to DB: %s", err.Error())
		return nil
	}
	log.Info("MongoDB => New client created")

	return client
}

func BuildTestMongoDatabase() *mongo.Database {
	clientOptions := new(options.ClientOptions)
	clientOptions.ApplyURI(ProjectSettings.Database.TestConnectionString)
	clientOptions.SetAppName(ProjectSettings.ProjectName)
	clientOptions.SetConnectTimeout(dbTimeout)
	clientOptions.SetSocketTimeout(dbSocketTimeout)
	clientOptions.SetServerSelectionTimeout(dbTimeout)
	clientOptions.SetMaxPoolSize(ProjectSettings.Database.MaxConnections)
	clientOptions.SetMaxConnIdleTime(dbMaxConnIdleTime)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Errorf("BuildMongoDatabaseClient | Error connecting to DB: %s", err.Error())
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), ctxTimeout)

	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("BuildMongoDatabaseClient | Error connecting to DB: %s", err.Error())
		return nil
	}
	log.Info("MongoDB => New client created")

	return client.Database(ProjectSettings.Database.TestMutantsDbName)
}