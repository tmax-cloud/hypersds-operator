package util

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
)

var _ = Describe("Ceph Exec Test", func() {
	defer GinkgoRecover()
	var (
		mockCtrl *gomock.Controller
		ioMock   *wrapper.MockIoUtilInterface
		execMock *wrapper.MockExecInterface
		osMock   *wrapper.MockOsInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		ioMock = wrapper.NewMockIoUtilInterface(mockCtrl)
		execMock = wrapper.NewMockExecInterface(mockCtrl)
		osMock = wrapper.NewMockOsInterface(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[RunCephCmd Test]", func() {
		It("Execute RunCephCmd", func() {
			ioMock.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			osMock.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			execMock.EXPECT().CommandExecute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			cephConf := []byte("[global]\nfsid = 1c7328c0-b6d9-11eb-bf40-b42e99095b95\nmon_host = [v2:192.168.7.19:3300/0,v1:192.168.7.19:6789/0]\n")
			cephKeyring := []byte("[client.admin]\nkey = AQAXDaJgnWT7LhAAVTXjbEKLN4bBMvVkWaqYLw==\n")
			cmd := []string{"-s"}
			// testBuf, err := RunCephCmd(wrapper.OsWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring, "test", cmd...)
			// fmt.Println(testBuf.String())
			_, err := RunCephCmd(osMock, execMock, ioMock, cephConf, cephKeyring, "test", cmd...)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
