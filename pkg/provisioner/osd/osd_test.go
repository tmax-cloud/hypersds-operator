package osd

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/service"
)

var _ = Describe("Osd Test", func() {
	defer GinkgoRecover()
	var (
		mockCtrl        *gomock.Controller
		yamlMock        *wrapper.MockYamlInterface
		yamlDecoderMock *wrapper.MockYamlDecoderInterface
		ioMock          *wrapper.MockIoUtilInterface
	)
	var (
		dataDevices Device
		s           service.Service
		placement   service.Placement
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		yamlMock = wrapper.NewMockYamlInterface(mockCtrl)
		ioMock = wrapper.NewMockIoUtilInterface(mockCtrl)
		yamlDecoderMock = wrapper.NewMockYamlDecoderInterface(mockCtrl)

		_ = placement.SetHosts([]string{"testHost"})
		_ = s.SetPlacement(placement)
		_ = s.SetServiceType("osd")
		_ = s.SetServiceID("osd_test")
		_ = dataDevices.setPaths([]string{"testPath", "testPath2"})
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Getter/Setter Test", func() {
		It("is simple case", func() {
			var osd Osd
			var err error

			err = osd.SetService(&s)
			Expect(err).NotTo(HaveOccurred())
			changedService := osd.GetService()
			Expect(changedService).To(Equal(s))

			err = osd.SetDataDevices(&dataDevices)
			Expect(err).NotTo(HaveOccurred())
			changedDataDevices := osd.GetDataDevices()
			Expect(changedDataDevices).To(Equal(dataDevices))
		})
	})
	Describe("MakeYmlFile Test", func() {

		It("should return nil", func() {
			var osd Osd
			_ = osd.SetService(&s)
			_ = osd.SetDataDevices(&dataDevices)
			fileName := "tmp.yaml"

			yamlMock.EXPECT().Marshal(gomock.Any()).Return(nil, nil).AnyTimes()
			ioMock.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			err := osd.MakeYmlFile(yamlMock, ioMock, fileName)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("NewOsdsFromCephCr Test", func() {

		It("is simple test case", func() {
			hostName := "testHost"
			devicePaths := []string{"testPath", "testPath2"}
			cephSpec := hypersdsv1alpha1.CephClusterSpec{
				Osd: []hypersdsv1alpha1.CephClusterOsdSpec{
					hypersdsv1alpha1.CephClusterOsdSpec{
						HostName: hostName,
						Devices:  devicePaths,
					},
				},
			}
			serviceType := "osd"
			serviceID := "osd_" + hostName

			osdList, err := NewOsdsFromCephCr(cephSpec)
			Expect(err).NotTo(HaveOccurred())
			osd := osdList[0]
			changedService := osd.GetService()
			changedPlacement := changedService.GetPlacement()
			changedDataDevices := osd.GetDataDevices()
			Expect(changedPlacement.GetHosts()[0]).To(Equal(hostName))
			Expect(changedService.GetServiceType()).To(Equal(serviceType))
			Expect(changedService.GetServiceID()).To(Equal(serviceID))
			Expect(changedDataDevices.getPaths()).To(Equal(devicePaths))
		})
	})
	Describe("NewOsdsFromCephOrch Test", func() {

		It("is simple test case", func() {
			yamlMock.EXPECT().NewDecoder(gomock.Any()).Return(yamlDecoderMock).AnyTimes()
			gomock.InOrder(
				yamlDecoderMock.EXPECT().Decode(gomock.Any()).DoAndReturn(func(osd *Osd) error {
					_ = osd.SetService(&s)
					_ = osd.SetDataDevices(&dataDevices)
					return nil
				}),
				yamlDecoderMock.EXPECT().Decode(gomock.Any()).Return(errors.New("error")),
			)
			rawOsdsFromOrch := []byte{}
			osdList, err := NewOsdsFromCephOrch(yamlMock, rawOsdsFromOrch)
			Expect(err).NotTo(HaveOccurred())
			osd := osdList[0]
			changedService := osd.GetService()
			changedDataDevices := osd.GetDataDevices()
			Expect(changedService).To(Equal(s))
			Expect(changedDataDevices).To(Equal(dataDevices))
		})
	})
	Describe("CompareDataDevices Test", func() {
		It("is simple test case", func() {
			var osd, targetOsd Osd
			var targetDataDevices Device

			_ = osd.SetService(&s)
			_ = osd.SetDataDevices(&dataDevices)

			_ = targetDataDevices.setPaths([]string{"testPath2", "addPath"})
			_ = targetOsd.SetService(&s)
			_ = targetOsd.SetDataDevices(&targetDataDevices)

			expectAddDeviceList := []string{"addPath"}
			expectRemoveDeviceList := []string{"testPath"}

			addDeviceList, removeDeviceList, err := osd.CompareDataDevices(&targetOsd)
			Expect(err).NotTo(HaveOccurred())
			Expect(addDeviceList).To(Equal(expectAddDeviceList))
			Expect(removeDeviceList).To(Equal(expectRemoveDeviceList))
		})
	})
})
