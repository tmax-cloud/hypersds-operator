package node

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OS Packager Test", func() {
	Context("strToPackager", func() {
		It("Should return error", func() {
			_, err := strToPackager("abc")
			Expect(err).To(HaveOccurred())
		})
		It("Should return apt for ubuntu", func() {
			packager, err := strToPackager("ubuntu")
			Expect(err).NotTo(HaveOccurred())
			Expect(packager).To(Equal(Apt))
		})
		It("Should return yum for centos", func() {
			packager, err := strToPackager("centos")
			Expect(err).NotTo(HaveOccurred())
			Expect(packager).To(Equal(Yum))
		})
	})
	Context("getAptDistros", func() {
		It("Should return distros using apt", func() {
			aptgets := []string{string(Ubuntu)}
			output := getAptDistros()
			Expect(output).To(Equal(aptgets))
		})
	})
	Context("getYumDistros", func() {
		It("Should return distros using yum", func() {
			yums := []string{string(Centos)}
			output := getYumDistros()
			Expect(output).To(Equal(yums))
		})
	})
	Context("contains", func() {
		It("Should not return error", func() {
			arr := []string{"abc"}
			Expect(contains(arr, "abc")).Should(BeTrue())
			Expect(contains(arr, "def")).Should(BeFalse())
		})
	})
})

var _ = Describe("OS Distro Test", func() {
	Context("strToDistro", func() {
		It("Should return error", func() {
			_, err := strToDistro("abc")
			Expect(err).To(HaveOccurred())
		})
		It("Should return ubuntu", func() {
			distro, err := strToDistro("ubuntu")
			Expect(err).NotTo(HaveOccurred())
			Expect(distro).To(Equal(Ubuntu))
		})
		It("Should return centos", func() {
			distro, err := strToDistro("centos")
			Expect(err).NotTo(HaveOccurred())
			Expect(distro).To(Equal(Centos))

		})
	})
})

var _ = Describe("OS Version Test", func() {
	Context("strToVersion", func() {
		It("Should return error", func() {
			_, err := strToVersion("abc")
			Expect(err).To(HaveOccurred())
		})
		It("Should parse decimal", func() {
			ver, err := strToVersion("20.04")
			Expect(err).NotTo(HaveOccurred())
			Expect(ver).To(Equal(20.04))
		})
		It("Should parse integer", func() {
			ver, err := strToVersion("20")
			Expect(err).NotTo(HaveOccurred())
			Expect(ver).To(Equal(float64(20)))
		})
	})
})
