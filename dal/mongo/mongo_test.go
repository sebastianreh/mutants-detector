package mongo_test

import (
	. "ExamenMeLiMutante/dal/mongo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("Mongo", func() {
	var mongoCollection = IMongoDatabase(MongoDatabase{})
	var mongoDatabase *mongo.Database

	Context("calling mongo server", func() {
		It("should return the redis server", func() {
			mongoCollection = NewMongoClient()
			mongoDatabase = mongoCollection.SubjectsDatabase()
			expected := mongo.Database{}
			Expect(mongoDatabase).Should(BeAssignableToTypeOf(&expected))
		})

		It("should return the mongo server", func() {
			mongoCollection = nil
			expected := mongo.Database{}
			Expect(mongoDatabase).Should(BeAssignableToTypeOf(&expected))
		})
	})

	Context("configuring database", func() {
		It("should configure database correctly", func() {
			mongoCollection = NewMongoClient()
			mongoCollection.ConfigureDatabase()
			mongoDatabase = mongoCollection.SubjectsDatabase()
			expected := mongo.Database{}
			Expect(mongoDatabase).Should(BeAssignableToTypeOf(&expected))
		})
	})
})
