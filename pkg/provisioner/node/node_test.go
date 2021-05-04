package node

import (
	"bytes"

	common "github.com/tmax-cloud/hypersds-operator/pkg/provisioner/common/wrapper"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"golang.org/x/crypto/ssh"
)

var _ = Describe("Node Test", func() {
	defer GinkgoRecover()

	var (
		testingNode                      Node
		userId, userPw, ipAddr, hostName string
		hostSpec                         HostSpec
		cephSpec                         hypersdsv1alpha1.CephClusterSpec
	)

	Describe("Getter/Setter Test", func() {
		It("is simple test case", func() {
			// userId getter/setter test
			userId = "shellwedance"
			err := testingNode.SetUserId(userId)
			Expect(err).NotTo(HaveOccurred())

			changedUserId := testingNode.GetUserId()
			Expect(changedUserId).To(Equal(userId))

			// userPw getter/setter test
			userPw = "123abc!@#"
			err = testingNode.SetUserPw(userPw)
			Expect(err).NotTo(HaveOccurred())

			changedUserPw := testingNode.GetUserPw()
			Expect(changedUserPw).To(Equal(userPw))

			// hostSpec getter/setter test
			hostSpec = HostSpec{
				ServiceType: HostSpecServiceType,
			}
			err = testingNode.SetHostSpec(hostSpec)
			Expect(err).NotTo(HaveOccurred())

			changedHostSpec := testingNode.GetHostSpec()
			Expect(changedHostSpec).To(Equal(hostSpec))
		})
	})

	Describe("RunSshCmd Test", func() {
		var (
			mockCtrl *gomock.Controller
			m        *common.MockSshInterface
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			m = common.NewMockSshInterface(mockCtrl)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("is simple test case", func() {
			testCommand := "hello world"
			m.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(addr, command string, resultStdout, resultStderr *bytes.Buffer, config *ssh.ClientConfig) error {
					resultStdout.WriteString("hello world")
					return nil
				}).AnyTimes()

			result, err := testingNode.RunSshCmd(m, testCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(testCommand))
		})
	})

	Describe("[NewNodesFromCephCr Test]", func() {
		It("is simple test case", func() {
			ipAddr = "1.1.1.1"
			userId = "developer1"
			userPw = "abc123!@#"
			hostName = "node1"

			cephSpec = hypersdsv1alpha1.CephClusterSpec{
				Nodes: []hypersdsv1alpha1.Node{
					{
						IP:       ipAddr,
						UserID:   userId,
						Password: userPw,
						HostName: hostName,
					},
				},
			}

			nodes, err := NewNodesFromCephCr(cephSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodes)).To(Equal(1))

			createdUserId := nodes[0].GetUserId()
			Expect(createdUserId).To(Equal(userId))

			createdUserPw := nodes[0].GetUserId()
			Expect(createdUserPw).To(Equal(userId))

			createdHostSpec := nodes[0].GetHostSpec()
			Expect(createdHostSpec.GetServiceType()).To(Equal(HostSpecServiceType))
			Expect(createdHostSpec.GetHostName()).To(Equal(hostName))
			Expect(createdHostSpec.GetAddr()).To(Equal(ipAddr))
			Expect(createdHostSpec.GetLabels()).To(BeEmpty())
			Expect(createdHostSpec.GetStatus()).To(BeEmpty())
		})
	})
})
