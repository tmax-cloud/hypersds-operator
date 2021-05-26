package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Placement Test", func() {
	defer GinkgoRecover()

	Describe("Getter/Setter Test", func() {
		It("is simple case", func() {
			var err error
			var placement Placement

			label := "testLabel"
			hosts := []string{"testHost"}
			count := 1
			hostPattern := "testHostPattern"

			err = placement.SetLabel(label)
			Expect(err).NotTo(HaveOccurred())
			changedLabel := placement.GetLabel()
			Expect(changedLabel).To(Equal(label))

			err = placement.SetHosts(hosts)
			Expect(err).NotTo(HaveOccurred())
			changedHosts := placement.GetHosts()
			Expect(changedHosts).To(Equal(hosts))

			err = placement.SetCount(count)
			Expect(err).NotTo(HaveOccurred())
			changedCount := placement.GetCount()
			Expect(changedCount).To(Equal(count))

			err = placement.SetHostPattern(hostPattern)
			Expect(err).NotTo(HaveOccurred())
			changedHostPattern := placement.GetHostPattern()
			Expect(changedHostPattern).To(Equal(hostPattern))

		})
	})
})
