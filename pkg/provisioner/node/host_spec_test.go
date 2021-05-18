package node

import (
	"os"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HostSpec Test", func() {
	defer GinkgoRecover()

	var (
		hostSpec                 HostSpec
		ipAddr, hostName, status string
		labels, labelsToAdd      []string
	)

	Describe("Getter/Setter Test", func() {
		// TODO: replace to ginkgo table extension
		It("is simple case", func() {
			// ServiceType getter/setter test
			err := hostSpec.SetServiceType()
			Expect(err).NotTo(HaveOccurred())

			createdServiceType := hostSpec.GetServiceType()
			Expect(createdServiceType).To(Equal(HostSpecServiceType))

			// Addr getter/setter test
			ipAddr = "1.1.1.2"
			err = hostSpec.SetAddr(ipAddr)
			Expect(err).NotTo(HaveOccurred())

			changedAddr := hostSpec.GetAddr()
			Expect(changedAddr).To(Equal(ipAddr))

			// HostName getter/setter test
			hostName = "node2"
			err = hostSpec.SetHostName(hostName)
			Expect(err).NotTo(HaveOccurred())

			changedHostName := hostSpec.GetHostName()
			Expect(changedHostName).To(Equal(hostName))

			// Labels getter/setter test
			labels = []string{"exampleA", "exampleB"}
			err = hostSpec.SetLabels(labels)
			Expect(err).NotTo(HaveOccurred())

			changedLabels := hostSpec.GetLabels()
			Expect(changedLabels).To(Equal(labels))

			// Labels adder test
			err = hostSpec.AddLabels(labelsToAdd...)
			Expect(err).NotTo(HaveOccurred())

			allLabels := hostSpec.GetLabels()
			labels = append(labels, labelsToAdd...)
			Expect(allLabels).To(Equal(labels))

			// HostName getter/setter test
			status = "CEPHADM_HOST_CHECK_FAILED"
			err = hostSpec.SetStatus(status)
			Expect(err).NotTo(HaveOccurred())

			changedStatus := hostSpec.GetStatus()
			Expect(changedStatus).To(Equal(status))
		})

	})

	Describe("MakeYmlFile Test", func() {
		var (
			mockCtrl *gomock.Controller
			mYaml    *wrapper.MockYamlInterface
			mIoUtil  *wrapper.MockIoUtilInterface
			fileName string
		)

		BeforeEach(func() {
			ipAddr = "1.1.1.1"
			hostName = "node1"
			status = "CEPHADM_STRAY_HOST"
			labels = []string{"example1", "example2"}
			labelsToAdd = []string{"example3", "example4"}

			hostSpec = HostSpec{
				ServiceType: HostSpecServiceType,
				Addr:        ipAddr,
				HostName:    hostName,
				Labels:      labels,
				Status:      status,
			}

			mockCtrl = gomock.NewController(GinkgoT())
			mYaml = wrapper.NewMockYamlInterface(mockCtrl)
			mIoUtil = wrapper.NewMockIoUtilInterface(mockCtrl)

			fileName = "tmp.yml"
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("is simple test case", func() {
			mYaml.EXPECT().Marshal(gomock.Any()).DoAndReturn(func(hs *HostSpec) ([]byte, error) {
				return nil, nil
			}).AnyTimes()
			mIoUtil.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(fileName string, data []byte, fileMode os.FileMode) error {
					return nil
				}).AnyTimes()

			err := hostSpec.MakeYmlFile(mYaml, mIoUtil, fileName)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
