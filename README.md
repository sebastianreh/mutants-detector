## Mutant-detector

Hi! This is the mutant detector Microservice

 ### Requirements

- Golang v1.13.x

- To run tests or using this microservice, please consider using docker, as you need both redis and MongoDB servers.

- To install dependecies:
`$ go mod download`

#### To run tests:
    
- Run docker dependencies:

  `docker run --rm -d --name redis-test -p 6379:6379 redis`

  `docker run --rm -d --name mongo-test -p 27017:27017 mongo`