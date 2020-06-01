package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MutantsCollection = "mutant_subjects"
	HumansCollection  = "humans_subjects"
)

var (
	unique = true
	sparse = true
)

func (mgo MongoDatabase) ConfigureDatabase() {
	log.Info("Configuring database...")
	db := mgo.SubjectsDatabase()
	if db != nil {
		ConfigureIndexes(db)
	} else {
		log.Error("Invalid database to configure")
	}
}

func ConfigureIndexes(db *mongo.Database) {
	setSubjectsIndexes(db)
}

func setSubjectsIndexes(db *mongo.Database) {
	setIndex(db, MutantsCollection)
	setIndex(db, HumansCollection)
}

func setIndex(db *mongo.Database, collectionName string) {
	collection := db.Collection(collectionName)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"dna_id": 1},
		Options:
		&options.IndexOptions{
			Unique: &unique,
			Sparse: &sparse,
		},
	})
	if err != nil {
		log.Error("dal.db_settings.setIndex | Error: ", err)
	}
}
