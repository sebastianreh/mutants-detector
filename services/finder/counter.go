package finder

import (
	"strings"
	"sync"
)

type (
	mutationsSafeCounter struct {
		mutations int64
		mux       sync.Mutex
	}
)

// Es un contador seguro para evitar los race conditions que puedan aparecer al contar las mutaciones

func NewMutationsSafeCounter() *mutationsSafeCounter{
	return &mutationsSafeCounter{}
}

// Cuenta las mutaciones que existen en cada string dentro de la cadena de and

func (counter *mutationsSafeCounter) findMutationsCounterAsync(dnaBit string) int64 {
	count := int64(strings.Count(dnaBit, mutationA) + strings.Count(dnaBit, mutationT) +
		strings.Count(dnaBit, mutationC) + strings.Count(dnaBit, mutationG))
	return count
}

// Bloquea e incrementa el contador si encuentro una mutación
// Cuando las mutaciones están por encima del threshold devuelve la condición de mutante y cierra el resto de las rutinas

func (counter *mutationsSafeCounter) increaseCountAndFind(dnaBit string, isMutant *chan bool, quitChan *chan struct{}) {
	numberOfMutations := counter.findMutationsCounterAsync(dnaBit)
	if numberOfMutations > 0 {
		counter.mux.Lock()
		counter.mutations += numberOfMutations
		if counter.mutations > mutationThreshold {
			*isMutant <- true
			close(*quitChan)
		}
		counter.mux.Unlock()
	}
}

// Cuando no se encuentran mutaciones, devuelve la condición de humano al canal

func (counter *mutationsSafeCounter) NoMutationFound(isMutant *chan bool) {
	counter.mux.Lock()
	*isMutant <- false
	counter.mux.Unlock()
}
