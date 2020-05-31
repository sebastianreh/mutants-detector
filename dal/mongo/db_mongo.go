package mongo

import (
	. "ExamenMeLiMutante/settings"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctxTimeout = time.Duration(30) * time.Second
var dbTimeout = time.Duration(ProjectSettings.Database.Timeout) * time.Second
var dbSocketTimeout = time.Duration(60) * time.Second
var dbMaxConnIdleTime = time.Duration(ProjectSettings.Database.MaxConnIdleTime) * time.Second

func BuildMongoDatabaseClient(connection string) *mongo.Client {
	clientOptions := new(options.ClientOptions)
	clientOptions.ApplyURI(connection)
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
