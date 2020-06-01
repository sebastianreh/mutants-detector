package repositories_test

import (
	"ExamenMeLiMutante/models"
	"ExamenMeLiMutante/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"syreclabs.com/go/faker"
)

var _ = Describe("MutantsCache", func() {

	var mutantsCache repositories.IMutantCache

	BeforeEach(func(){
		mutantsCache = repositories.NewMutantCacheRepository()
	})

	Context("when saving stats in cache", func() {
		It("should save in cache", func() {
			stats := models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt64(1, 100)),
				CountHumanDna:  float64(faker.RandomInt64(1, 100)),
			}
			stats.Ratio = stats.CountMutantDna/stats.CountHumanDna
			mutantsCache.SaveStatsInCache(&stats)
		})
	})

	Context("when getting stats from cache", func() {
		It("should return nil stats from cache", func() {
			stats := mutantsCache.GetStatsFromCache()
			Expect(stats).Should(BeNil())
		})

		It("should return stats from cache", func() {
			stats := models.MutantsStats{
				CountMutantDna: float64(faker.RandomInt64(1, 100)),
				CountHumanDna:  float64(faker.RandomInt64(1, 100)),
			}
			stats.Ratio = stats.CountMutantDna/stats.CountHumanDna
			mutantsCache.SaveStatsInCache(&stats)
			savedStats := mutantsCache.GetStatsFromCache()
			Expect(stats.CountMutantDna).Should(Equal(savedStats.CountMutantDna))
			Expect(stats.CountHumanDna).Should(Equal(savedStats.CountHumanDna))
			Expect(stats.Ratio).Should(Equal(savedStats.Ratio))
		})
	})
})
