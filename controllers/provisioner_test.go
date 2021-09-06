package controllers

//	no.	configmap	secret
//	1	x	x
//	2	o	x
//	3	x	o

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("syncProvisioner", func() {
	Context("1. without any resources", func() {
		r := createFakeCephClusterReconciler()
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("2. without secret", func() {
		cm := newConfigMap()
		r := createFakeCephClusterReconciler(cm)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("3. without configmap", func() {
		s := newSecret()
		r := createFakeCephClusterReconciler(s)
		err := r.syncProvisioner()

		It("Should not return error", func() {
			Expect(err).Should(BeNil())
		})
	})
})
