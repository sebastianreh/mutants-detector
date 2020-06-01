package redis_test

import (
	."github.com/sebastianreh/mutants-detector/dal/redis"
	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redis", func() {
	var redisDatabase = IRedisDatabase(RedisDatabase{})
	var redisClient *redis.Client

	Context("calling redis server", func() {
		It("should return the redis server", func() {
			redisDatabase = NewRedisDatabase()
			redisClient = redisDatabase.RedisClient()
			expected := redis.Client{}
			Expect(redisClient).Should(BeAssignableToTypeOf(&expected))
		})

		It("should return a new redis server", func() {
			redisClient = nil
			expected := redis.Client{}
			Expect(redisClient).Should(BeAssignableToTypeOf(&expected))
		})
	})
})
