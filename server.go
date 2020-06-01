package main

import (
	. "github.com/sebastianreh/mutants-detector/container"
	"github.com/sebastianreh/mutants-detector/server"
)

func main() {
	go MongoDatabase.ConfigureDatabase()
	go MutantRepository.BuildRedisData()
	server.SetupServer()
}
