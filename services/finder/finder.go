package finder

import (
	. "sync"
)

type (
	MutantFinderService struct {
	}

	IMutantFinderService interface {
		IsMutant(dna []string) bool
	}
)

const (
	mutationA           = "AAAA"
	mutationT           = "TTTT"
	mutationC           = "CCCC"
	mutationG           = "GGGG"
	mutationChainLength = 4
	mutationThreshold   = 1
)

var (
	offset = mutationChainLength
)

func NewMutantFinderService() IMutantFinderService {
	return MutantFinderService{}
}

	// Aqui comienza la busqueda de mutantes. Se crea un canal que espera por la respuesta de la condición del sujeto

func (finder MutantFinderService) IsMutant(dna []string) bool {
	isMutant := make(chan bool)
	sync := NewWaitGroupSyncronizer()
	go finder.findAsync(&dna, sync, &isMutant)
	return <-isMutant
}

	// Se lanza una rutina por cada tipo de busqueda, por adn: horizontal, vertical, diagonal y la diagonal inversa

func (finder MutantFinderService) findAsync(dnaChain *[]string, sync *waitGroupSynchronizer, isMutant *chan bool) {
	mutationCounter := NewMutationsSafeCounter()
	sync.general.Add(4)
	go finder.findHorizontalAsync(dnaChain, mutationCounter, sync, isMutant)
	go finder.findVerticalAsync(dnaChain, mutationCounter, sync, isMutant)
	go finder.findDiagonalAsync(dnaChain, mutationCounter, sync, isMutant)
	go finder.findInverseDiagonalAsync(dnaChain, mutationCounter, sync, isMutant)
	sync.general.Wait()
	mutationCounter.NoMutationFound(isMutant)
}

	// Lanza una rutina por cada string dentro del adn

func (finder MutantFinderService) findHorizontalAsync(dnaChain *[]string, mutations *mutationsSafeCounter, sync *waitGroupSynchronizer, isMutant *chan bool) {
	for _, dnaBit := range *dnaChain {
		sync.horizontal.Add(1)
		go finder.findMutationRoutine(mutations, dnaBit, &sync.horizontal, isMutant)
	}
	sync.horizontal.Wait()
	sync.general.Done()
}

	// Se arman las verticales una a una y se lanzan rutinas para su busqueda

func (finder MutantFinderService) findVerticalAsync(dnaChain *[]string, mutations *mutationsSafeCounter, sync *waitGroupSynchronizer, isMutant *chan bool) {
	for i := 0; i < len(*dnaChain); i++ {
		var tempDnaBit string
		for j := 0; j < len((*dnaChain)[i]); j++ {
			tempDnaBit = tempDnaBit + string((*dnaChain)[j][i])
		}
		sync.vertical.Add(1)
		go finder.findMutationRoutine(mutations, tempDnaBit, &sync.vertical, isMutant)
	}
	sync.vertical.Wait()
	sync.general.Done()
}

	// Se arman las diagonales, excluyendo los valores en las esquinas de la matriz que no pueden ser interpretados, por la longitud de la combinación
	// que se está buscando

func (finder MutantFinderService) findDiagonalAsync(dnaChain *[]string, mutations *mutationsSafeCounter, sync *waitGroupSynchronizer, isMutant *chan bool) {
	searchLimit := len(*dnaChain) - offset + 1
	for i := 0; i < searchLimit; i++ {
		var tempDnaBit string
		k := offset + i - 1
		for j := 0; j < offset+i; j++ {
			tempDnaBit = tempDnaBit + string((*dnaChain)[k][j])
			k--
		}
		sync.diagonal.Add(1)
		go finder.findMutationRoutine(mutations, tempDnaBit, &sync.diagonal, isMutant)
	}
	searchLimit = len(*dnaChain) - offset
	for i := 0; i < searchLimit; i++ {
		var tempDnaBit string
		k := len(*dnaChain) - 1
		for j := 1 + i; j < len(*dnaChain); j++ {
			tempDnaBit = tempDnaBit + string((*dnaChain)[k][j])
			k--
		}
		sync.diagonal.Add(1)
		go finder.findMutationRoutine(mutations, tempDnaBit, &sync.diagonal, isMutant)
	}
	sync.diagonal.Wait()
	sync.general.Done()
}

	// Se arman las diagonales, excluyendo los valores en las esquinas de la matriz que no pueden ser interpretados, por la longitud de la combinación
	// que se está buscando

func (finder MutantFinderService) findInverseDiagonalAsync(dnaChain *[]string, mutations *mutationsSafeCounter, sync *waitGroupSynchronizer, isMutant *chan bool) {
	searchLimit := len(*dnaChain) - offset + 1
	for i := 0; i < searchLimit; i++ {
		var tempDnaBit string
		k := offset + i - 1
		for j := len(*dnaChain) - 1; j > len(*dnaChain)-offset-i-1; j-- {
			tempDnaBit = tempDnaBit + string((*dnaChain)[k][j])
			k--
		}
		sync.inverseDiagonal.Add(1)
		go finder.findMutationRoutine(mutations, tempDnaBit, &sync.inverseDiagonal, isMutant)
	}
	searchLimit = len(*dnaChain) - offset
	for i := 0; i < searchLimit; i++ {
		var tempDnaBit string
		k := len(*dnaChain) - 1
		for j := len(*dnaChain) - 2 - i; j > -1; j-- {
			tempDnaBit = tempDnaBit + string((*dnaChain)[k][j])
			k--
		}
		sync.inverseDiagonal.Add(1)
		go finder.findMutationRoutine(mutations, tempDnaBit, &sync.inverseDiagonal, isMutant)
	}
	sync.inverseDiagonal.Wait()
	sync.general.Done()
}

	// Cuando recibe el evento de cierre del canal, deja de enviar nuevas rutinas para encontrar la mutación y aumentar el contador

func (finder MutantFinderService) findMutationRoutine(mutations *mutationsSafeCounter, dnaBit string, wg *WaitGroup, isMutant *chan bool) {
	var QuitChan = make(chan struct{})
	defer wg.Done()
	select {
	case <-QuitChan:
		return
	default:
		mutations.increaseCountAndFind(dnaBit, isMutant, &QuitChan)
	}
}
