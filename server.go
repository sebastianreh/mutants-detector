package main

import (
	. "ExamenMeLiMutante/container"
	"ExamenMeLiMutante/server"
)

func main() {
	go MongoDatabase.ConfigureDatabase()
	go MutantRepository.BuildRedisData()
	server.SetupServer()
}
