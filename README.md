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

#### How does it work?

When a mutant detection request is made, first, it validates that the requirements are fulfilled: that means, that the request structure is correct, that there are no other characters than those which are valid and that the NxN structure is also respected. 
After the validation is made, the search algorithm starts to find possible mutations. To acomplish this process more rapidly, it fires 4 goroutines, and each one fires one more routine per string in the array of strings. For each mutation found, a counter object (assigned to that dna) with a mutex (to avoid possible reace conditions) adds one to its counter. When the number of mutations are above the limit (in this case 2), the counter send the confirmation that the dna is mutant via the "isMutant" bool channel that waits for this value. If all the goroutines are completed and the counter has not raised above the mutations limit, it means that no mutations were found, so the "isMutant" channel sends a false value. After a value is sent, the channel gets closed.
When the service finishes to process the mutations, the store strategy begins. In order to support a really big number of requests per second, the project, uses Redis for a rapid KVstorage and MongoDB as his main DB. The process goes like this: All the DNAs that are processed by this service, get converted into a unique id that has this structure: The first character  of the id is a number that represents the dimensions of the NxN dna, and the following characters are the all its dna strings. This way, the dna can be easily converted from a structure to an ID and backwards. After this process, the ID is searched in Redis to see if it was already processed and/or stored. If it was not, the service processes this dna and saves it in a notSaved key. If the service gets more than X requests or Y seconds have passed after a previous requets (in this case 1000 requests and 5 seconds are configured), the service fires the permanently save process.
All the Ids that are in the not saved key are compared to those in the saved key (although they was also a previous confirmation to see if they were repeated) and all the new and unique IDs get converted into a DNA object and bulk written into the database. There is one collection for mutants and one for humans. The IDs are the string indexes of the collections, and all the subjects that are requested to be saved into the DB are ensured to be unique so the process is really fast. At the same time, all the IDs that were saved into the DB are now in the saved key, to also ensure that they are now permanently saved into the DB. 
When a stats request its made the service can easilly and effortlessly read the counter of each redis key, calculate the ratio, and then save its value into a cache. When a new dna is saved, the cache gets reset and the stats are once again calculated and stored into the cache. 

