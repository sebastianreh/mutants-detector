package services

import (
	"github.com/sebastianreh/mutants-detector/models"
	"github.com/sebastianreh/mutants-detector/repositories"
	"github.com/sebastianreh/mutants-detector/services/finder"
	"github.com/sebastianreh/mutants-detector/utils"
	log "github.com/sirupsen/logrus"
	"sync"
)

type (
	MutantService struct {
		finder     finder.IMutantFinderService
		repository repositories.IMutantRepository
	}

	IMutantService interface {
		VerifyMutant(mutantRequest models.MutantRequest) models.MutantResponse
		GetSubjectsStats() (*models.MutantsStats, error)
		ChangeCacheStatus(status bool)
	}
)

var (
	statsUpdated bool
	subChan      = make(chan *models.Subject)
	wg           sync.WaitGroup
)

func NewMutantService(finderService finder.IMutantFinderService, repository repositories.IMutantRepository) IMutantService {
	return MutantService{
		finder:     finderService,
		repository: repository,
	}
}

// Aqui se verifica el status del sujeto, si no está guardado se agrega se envia al canal
// y comienza la cola de guaradado

func (service MutantService) VerifyMutant(mutantRequest models.MutantRequest) models.MutantResponse {
	dnaId := utils.ConvertDnaToId(mutantRequest.DnaChain)
	subjectStatus := service.repository.GetSubjectStatus(dnaId)
	subject := new(models.Subject)
	*subject = models.Subject{
		Id:  dnaId,
		Dna: mutantRequest.DnaChain,
	}
	switch subjectStatus {
	case repositories.MutantStatus:
		subject.IsMutant = true
	case repositories.HumanStatus:
		subject.IsMutant = false
	default:
		subject.IsMutant = service.finder.IsMutant(mutantRequest.DnaChain)
		service.ChangeCacheStatus(false)
		go service.repository.QueueDatabaseOperations(subChan)
		subChan <- subject
	}
	return models.MutantResponse{IsMutant: subject.IsMutant}
}

// Aquí se consultan las stats, las mismas son calculadas o traidas del cache.

func (service MutantService) GetSubjectsStats() (*models.MutantsStats, error) {
	var stats *models.MutantsStats
	if !statsUpdated {
		preStats, err := service.repository.GetSubjectsStats()
		if err != nil {
			log.Errorf("services.mutant.GetSubjectStats error | %v", err)
			return nil, err
		}
		stats = utils.CalculateMutantStats(preStats)
		service.repository.SaveStatsInCache(stats)
		service.ChangeCacheStatus(true)
	} else {
		stats = service.repository.GetStatsFromCache()
	}
	log.Infof("services.mutant.GetSubjectStats | Stats requested: Mutants: %d Humans: %d Ratio: %f", int(stats.CountMutantDna), int(stats.CountHumanDna), stats.Ratio)
	return stats, nil
}

func (service MutantService) ChangeCacheStatus(status bool) {
	statsUpdated = status
}
