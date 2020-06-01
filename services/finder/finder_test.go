package finder_test

import (
	"ExamenMeLiMutante/services/finder"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {

	var ctrl *gomock.Controller
	var dna []string

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when a subject request needs to be searched", func() {
		It("should return mutant true", func() {
			dna = []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
			finderService := finder.NewMutantFinderService()
			isMutat := finderService.IsMutant(dna)
			Expect(isMutat).To(BeTrue())
		})

		It("should return mutant false", func() {
			dna = []string{"ATGC", "CGTA", "ATGC", "CGTG"}
			finderService := finder.NewMutantFinderService()
			isMutat := finderService.IsMutant(dna)
			Expect(isMutat).To(BeFalse())
		})
	})
})