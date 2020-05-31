package dal

import (
	. "ExamenMeLiMutante/settings"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type (
	RedisDatabase struct {
		client *redis.Client
	}

	IRedisDatabase interface {
		RedisClient() *redis.Client
	}
)

func NewRedisDatabase() IRedisDatabase {
	return RedisDatabase{
	}
}

func (rs RedisDatabase) RedisClient() *redis.Client {
	if rs.client != nil {
		return rs.client
	}
	rs.client = rs.buildClient(ProjectSettings.Redis.RedisHost, ProjectSettings.Redis.RedisPort)
	return rs.client
}

func (rs RedisDatabase) buildClient(host string, port string) *redis.Client {
	var url = fmt.Sprintf("%s:%s", host, port)
	var options = &redis.Options{
		Addr:     url,
		Password: "", // No password set
		DB:       0,  // Use default DB
		OnConnect: func(*redis.Conn) error {
			return nil
		},
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolSize:     100,
		PoolTimeout:  30 * time.Second,
	}

	return redis.NewClient(options)
}
