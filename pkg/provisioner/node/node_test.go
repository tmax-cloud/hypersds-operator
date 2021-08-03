package node

import (
	"bytes"
	"context"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"golang.org/x/crypto/ssh"
)

var _ = Describe("Node Test", func() {
	defer GinkgoRecover()

	var (
		testingNode                      Node
		userID, userPw, ipAddr, hostName string
		hostSpec                         HostSpec
		cephSpec                         hypersdsv1alpha1.CephClusterSpec
	)

	Describe("Getter/Setter Test", func() {
		It("is simple test case", func() {
			// userID getter/setter test
			userID = "shellwedance"
			err := testingNode.SetUserID(userID)
			Expect(err).NotTo(HaveOccurred())

			changedUserID := testingNode.GetUserID()
			Expect(changedUserID).To(Equal(userID))

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
			err = testingNode.SetHostSpec(&hostSpec)
			Expect(err).NotTo(HaveOccurred())

			changedHostSpec := testingNode.GetHostSpec()
			Expect(changedHostSpec).To(Equal(hostSpec))

			// os distro getter/setter test
			err = testingNode.SetOs(Ubuntu, Apt, 20.04)
			Expect(err).NotTo(HaveOccurred())

			updated := testingNode.GetOs()
			Expect(updated.Distro).To(Equal(Ubuntu))
			Expect(updated.Packager).To(Equal(Apt))
			Expect(updated.Version).To(Equal(20.04))
		})
	})

	Describe("RunSSHCmd Test", func() {
		var (
			mockCtrl *gomock.Controller
			m        *wrapper.MockSSHInterface
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			m = wrapper.NewMockSSHInterface(mockCtrl)
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

			result, err := testingNode.RunSSHCmd(m, testCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(testCommand))
		})
	})

	Describe("RunScpCmd Test", func() {
		var (
			mockCtrl *gomock.Controller
			m        *wrapper.MockExecInterface
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			m = wrapper.NewMockExecInterface(mockCtrl)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("is simple test case", func() {
			testString := "hello world"
			role := DESTINATION
			srcFile := "src"
			destFile := "dest"
			m.EXPECT().CommandExecute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, resultStdout, resultStderr *bytes.Buffer, name string, arg ...string) error {
					resultStdout.WriteString(testString)
					return nil
				}).AnyTimes()

			result, err := testingNode.RunScpCmd(m, srcFile, destFile, role)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(testString))

			role = SOURCE
			result, err = testingNode.RunScpCmd(m, srcFile, destFile, role)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(testString))
		})
	})

	Describe("[NewNodesFromCephCr Test]", func() {
		It("is simple test case", func() {
			ipAddr = "1.1.1.1"
			userID = "developer1"
			userPw = "abc123!@#"
			hostName = "node1"

			cephSpec = hypersdsv1alpha1.CephClusterSpec{
				Nodes: []hypersdsv1alpha1.Node{
					{
						IP:       ipAddr,
						UserID:   userID,
						Password: userPw,
						HostName: hostName,
					},
				},
			}

			nodes, err := NewNodesFromCephCr(cephSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodes)).To(Equal(1))

			createdUserID := nodes[0].GetUserID()
			Expect(createdUserID).To(Equal(userID))

			createdUserPw := nodes[0].GetUserID()
			Expect(createdUserPw).To(Equal(userID))

			createdHostSpec := nodes[0].GetHostSpec()
			Expect(createdHostSpec.GetServiceType()).To(Equal(HostSpecServiceType))
			Expect(createdHostSpec.GetHostName()).To(Equal(hostName))
			Expect(createdHostSpec.GetAddr()).To(Equal(ipAddr))
			Expect(createdHostSpec.GetLabels()).To(BeEmpty())
			Expect(createdHostSpec.GetStatus()).To(BeEmpty())
		})
	})
})
