package services_test

import (
	"github.com/sebastianreh/mutants-detector/models"
	"github.com/sebastianreh/mutants-detector/repositories"
	"github.com/sebastianreh/mutants-detector/services"
	repositoryMock "github.com/sebastianreh/mutants-detector/test/mocks/repositories"
	finderMock "github.com/sebastianreh/mutants-detector/test/mocks/services/finder"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"syreclabs.com/go/faker"
)

var _ = Describe("Services", func() {
	var ctrl *gomock.Controller
	var mockFinderService *finderMock.MockIMutantFinderService
	var mockRepository *repositoryMock.MockIMutantRepository
	var mutantService services.IMutantService
	var mutantsReq models.MutantRequest
	var mutantsRes models.MutantResponse
	var mutantStats models.MutantsStats
	var mutantsPreStats models.MutantsPreStats

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockFinderService = finderMock.NewMockIMutantFinderService(ctrl)
		mockRepository = repositoryMock.NewMockIMutantRepository(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when a mutant request is already in database", func() {
		It("should return mutant ok", func() {
			dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}

			mutantsReq = models.MutantRequest{DnaChain: dna}
			mockRepository.EXPECT().GetSubjectStatus(gomock.Any()).Return(repositories.MutantStatus).AnyTimes()
			mutantService = services.NewMutantService(mockFinderService, mockRepository)
			mutantsRes = mutantService.VerifyMutant(mutantsReq)
			expectedMutantsRes := models.MutantResponse{IsMutant: true}

			Expect(mutantsRes).Should(Equal(expectedMutantsRes))
		})

		It("should return human ok", func() {
			dna := []string{"ATGC", "CGTA", "ATGC", "CGTG"}
			mutantsReq = models.MutantRequest{DnaChain: dna}
			mockRepository.EXPECT().GetSubjectStatus(gomock.Any()).Return(repositories.HumanStatus).AnyTimes()
			mutantService = services.NewMutantService(mockFinderService, mockRepository)
			mutantsRes = mutantService.VerifyMutant(mutantsReq)
			expectedMutantsRes := models.MutantResponse{IsMutant: false}

			Expect(mutantsRes).Should(Equal(expectedMutantsRes))
		})
	})
	Context("when a stats are requested", func() {
		It("should return new mutant stats calculated", func() {
			mutantStats = models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt(100, 1000)),
				CountHumanDna:  float64(faker.RandomInt(1000, 10000)),
			}
			mutantStats.Ratio = mutantStats.CountMutantDna / mutantStats.CountHumanDna
			mutantsPreStats = models.MutantsPreStats{
				CountMutants: int(mutantStats.CountMutantDna),
				CountHumans:  int(mutantStats.CountHumanDna),
			}
			mockRepository.EXPECT().GetSubjectsStats().Return(mutantsPreStats, nil).AnyTimes()
			mockRepository.EXPECT().SaveStatsInCache(gomock.Any())
			mutantService = services.NewMutantService(mockFinderService, mockRepository)
			expectedMutantStats, _ := mutantService.GetSubjectsStats()
			Expect(mutantStats).Should(Equal(*expectedMutantStats))
		})

		It("should return new mutant stats calculated", func() {
			mutantStats = models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt(100, 1000)),
				CountHumanDna:  float64(faker.RandomInt(1000, 10000)),
			}
			mutantStats.Ratio = mutantStats.CountMutantDna / mutantStats.CountHumanDna
			mutantsPreStats = models.MutantsPreStats{
				CountMutants: int(mutantStats.CountMutantDna),
				CountHumans:  int(mutantStats.CountHumanDna),
			}
			mutantService.ChangeCacheStatus(false)
			err := fmt.Errorf("services.mutant.GetSubjectStats error")
			mockRepository.EXPECT().GetSubjectsStats().Return(mutantsPreStats, err).AnyTimes()
			mutantService = services.NewMutantService(mockFinderService, mockRepository)
			mutantStats, err := mutantService.GetSubjectsStats()
			Expect(mutantStats).Should(BeNil())
			Expect(err).To(Equal(err))
		})

		It("should return mutant stats from cache", func() {
			mutantStats = models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt(100, 1000)),
				CountHumanDna:  float64(faker.RandomInt(1000, 10000)),
			}
			mutantStats.Ratio = mutantStats.CountMutantDna / mutantStats.CountHumanDna
			mockRepository.EXPECT().GetStatsFromCache().Return(&mutantStats).AnyTimes()
			mutantService = services.NewMutantService(mockFinderService, mockRepository)
			mutantService.ChangeCacheStatus(true)
			expectedMutantStats, _ := mutantService.GetSubjectsStats()
			Expect(mutantStats).Should(Equal(*expectedMutantStats))
		})
	})
})
