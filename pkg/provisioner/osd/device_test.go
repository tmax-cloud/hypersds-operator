package osd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Device Test", func() {
	defer GinkgoRecover()

	Describe("Getter/Setter Test", func() {
		It("is simple case", func() {
			var device Device
			var err error

			paths := []string{"testPath", "testPath2"}
			model := "testModel"
			size := "testSize"
			rotational := true
			vendor := "testVendor"
			all := false
			limit := 1

			err = device.setPaths(paths)
			Expect(err).NotTo(HaveOccurred())
			changedPaths := device.getPaths()
			Expect(changedPaths).To(Equal(paths))

			err = device.setModel(model)
			Expect(err).NotTo(HaveOccurred())
			changedModel := device.getModel()
			Expect(changedModel).To(Equal(model))

			err = device.setSize(size)
			Expect(err).NotTo(HaveOccurred())
			changedSize := device.getSize()
			Expect(changedSize).To(Equal(size))

			err = device.setRotational(rotational)
			Expect(err).NotTo(HaveOccurred())
			changedRotational := device.getRotational()
			Expect(changedRotational).To(Equal(rotational))

			err = device.setVendor(vendor)
			Expect(err).NotTo(HaveOccurred())
			changedVendor := device.getVendor()
			Expect(changedVendor).To(Equal(vendor))

			err = device.setAll(all)
			Expect(err).NotTo(HaveOccurred())
			changedAll := device.getAll()
			Expect(changedAll).To(Equal(all))

			err = device.setLimit(limit)
			Expect(err).NotTo(HaveOccurred())
			changedLimit := device.getLimit()
			Expect(changedLimit).To(Equal(limit))
		})
	})
})
