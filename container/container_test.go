package container_test

import (
	. "ExamenMeLiMutante/container"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	Context("when starting dependencies", func() {
		It("should return the dependencies", func() {
			Expect(MutantController).To(BeEquivalentTo(MutantController))
			Expect(MutantService).To(BeEquivalentTo(MutantService))
			Expect(MutantFinderService).To(BeEquivalentTo(MutantFinderService))
			Expect(MutantRepository).To(BeEquivalentTo(MutantRepository))
			Expect(MutantCache).To(BeEquivalentTo(MutantCache))
			Expect(MongoDatabase).To(BeEquivalentTo(MongoDatabase))
			Expect(RedisDatabase).To(BeEquivalentTo(RedisDatabase))
		})
	})
})
