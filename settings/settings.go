package settings

import "github.com/kelseyhightower/envconfig"

type (
	Settings struct {
		ProjectName    string `default:"mutant-finder"`
		ProjectVersion string `default:"0.1.0"`
		UrlBase        string `envconfig:"BASE_URL" required:"true" default:"localhost"`
		Host           string `envconfig:"HOST" default:"0.0.0.0"`
		Port           string `envconfig:"PORT" default:"8080"`
		Database       DatabaseSpecification
		Redis          RedisSpecification
	}

	DatabaseSpecification struct {
		ConnectionString     string `envconfig:"DB_MUTANTS_CONNECTION" required:"true"`
		MutantsDbName        string `envconfig:"DB_MUTANTS_COLLECTION" default:"mutant_subjects"`
		MaxConnections       uint64 `envconfig:"DB_MAX_CONNECTIONS" default:"20"`
		Timeout              int    `envconfig:"DB_TIMEOUT" default:"10"`
		MaxConnIdleTime      int    `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"30"`
		TestConnectionString string `default:"mongodb://localhost:27017/test-db"`
		TestMutantsDbName    string `default:"test-collection"`
	}

	RedisSpecification struct {
		RedisHost     string `envconfig:"REDIS_HOST" required:"true"`
		RedisPort     string `envconfig:"REDIS_PORT" required:"true"`
		TestRedisHost string `default:"127.0.0.1"`
	}
)

var ProjectSettings Settings

// Inicializa los settings del proyecyo. Si falta alguna variable de entorno lanza un panic.

func init() {
	ProjectSettings = Settings{}
	if err := envconfig.Process("", &ProjectSettings); err != nil {
		panic(err.Error())
	}
}
