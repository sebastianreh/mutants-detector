package finder

import "sync"

type (
	waitGroupSynchronizer struct {
		general         sync.WaitGroup
		horizontal      sync.WaitGroup
		vertical        sync.WaitGroup
		diagonal        sync.WaitGroup
		inverseDiagonal sync.WaitGroup
	}
)

// Crea un wait Group para sincronizar el buscado asincrónico de mutaciones

func NewWaitGroupSyncronizer() *waitGroupSynchronizer {
	return new(waitGroupSynchronizer)
}
