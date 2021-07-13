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
			var s Service
			var placement Placement

			serviceType := "testType"
			serviceID := "testId"
			unmanaged := false
			_ = placement.SetHosts([]string{"testHost"})

			err = s.SetServiceType(serviceType)
			Expect(err).NotTo(HaveOccurred())
			changedServiceType := s.GetServiceType()
			Expect(changedServiceType).To(Equal(serviceType))

			err = s.SetServiceID(serviceID)
			Expect(err).NotTo(HaveOccurred())
			changedServiceID := s.GetServiceID()
			Expect(changedServiceID).To(Equal(serviceID))

			err = s.SetPlacement(placement)
			Expect(err).NotTo(HaveOccurred())
			changedPlacement := s.GetPlacement()
			Expect(changedPlacement).To(Equal(placement))

			err = s.SetUnmanaged(unmanaged)
			Expect(err).NotTo(HaveOccurred())
			changedUnmanaged := s.GetUnmanaged()
			Expect(changedUnmanaged).To(Equal(unmanaged))

		})
	})
})
