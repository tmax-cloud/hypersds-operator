package config

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Config Test", func() {
	defer GinkgoRecover()

	var (
		mockCtrl *gomock.Controller
		ioMock   *wrapper.MockIoUtilInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		ioMock = wrapper.NewMockIoUtilInterface(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Getter/Setter Test", func() {
		It("is simple test case", func() {
			var testConfig CephConfig

			testConf := map[string]string{
				"fsid":     "b29fd",
				"mon_host": "[0.0.0.0]",
			}
			testSecret := map[string][]byte{
				"keyring": []byte("[client.admin]\n\tkey = b29fd"),
			}

			// crConf getter/setter test
			err := testConfig.SetCrConf(testConf)
			Expect(err).NotTo(HaveOccurred())

			changedConf := testConfig.GetCrConf()
			Expect(changedConf).To(Equal(testConf))

			// admConf getter/setter test
			err = testConfig.SetAdmConf(testConf)
			Expect(err).NotTo(HaveOccurred())

			changedConf = testConfig.GetAdmConf()
			Expect(changedConf).To(Equal(testConf))

			// admSecret getter/setter test
			err = testConfig.SetAdmSecret(testSecret)
			Expect(err).NotTo(HaveOccurred())

			changedSecret := testConfig.GetAdmSecret()
			Expect(changedSecret).To(Equal(testSecret))
		})
	})

	Describe("[NewNodesFromCephCr Test]", func() {
		It("is simple test case", func() {
			testConf := map[string]string{
				"fsid":     "b29fd",
				"mon_host": "[0.0.0.0]",
			}
			cephSpec := hypersdsv1alpha1.CephClusterSpec{
				Config: testConf,
			}

			config, err := NewConfigFromCephCr(cephSpec)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.GetCrConf()).To(Equal(testConf))
		})
	})

	Describe("[ConfigFromAdm Test]", func() {
		It("Parse ceph.conf to AdmConfig", func() {
			ioMock.EXPECT().ReadFile(gomock.Any()).DoAndReturn(
				func(filename string) ([]byte, error) {
					conf := []byte("[global]\n\tfsid = b29fd\n\tmon_host = [0.0.0.0]\n")
					return conf, nil
				}).AnyTimes()
			testConfig := CephConfig{}
			AdmConfig := map[string]string{
				"conf":     "[global]\n\tfsid = b29fd\n\tmon_host = [0.0.0.0]\n",
				"fsid":     "b29fd",
				"mon_host": "[0.0.0.0]",
			}
			err := testConfig.ConfigFromAdm(ioMock, "ceph.conf")
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfig.GetAdmConf()).To(Equal(AdmConfig))
		})
	})

	Describe("[SecretFromAdm Test]", func() {
		It("Parse keyring to AdmSecret", func() {
			ioMock.EXPECT().ReadFile(gomock.Any()).DoAndReturn(
				func(filename string) ([]byte, error) {
					secret := []byte("[client.admin]\n\tkey = b29fd")
					return secret, nil
				}).AnyTimes()
			testConfig := CephConfig{}
			AdmSecret := map[string][]byte{
				"keyring": []byte("[client.admin]\n\tkey = b29fd"),
			}
			err := testConfig.SecretFromAdm(ioMock, "keyring")
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfig.GetAdmSecret()).To(Equal(AdmSecret))
		})
	})

	Describe("[MakeIniFile Test]", func() {
		It("Make Ini file from Map", func() {
			ioMock.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(fileName string, data []byte, fileMode os.FileMode) error {
					return nil
				}).AnyTimes()

			testConfig := CephConfig{
				crConf: map[string]string{
					"debug_osd": "20/20",
				},
			}
			//ini := "[global]\n\tdebug_osd = 20/20\n"
			err := testConfig.MakeIniFile(ioMock, "ceph.conf")
			Expect(err).NotTo(HaveOccurred())
			//Expect(retini).To(Equal(ini))
		})
	})
	Describe("[INCOMPLETE][UpdateConfToK8s Test]", func() {
		testConfig := CephConfig{
			crConf: map[string]string{
				"test": "test",
			},
		}
		It("should return nil with configmap", func() {
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
			}
			fakeClient := fake.NewFakeClient(configMap)

			err := testConfig.UpdateConfToK8s(fakeClient, "test", "test")
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return err with no configmap", func() {
			fakeClient := fake.NewFakeClient()
			err := testConfig.UpdateConfToK8s(fakeClient, "test", "test")
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("[INCOMPLETE][UpdateKeyringToK8s Test]", func() {
		testConfig := CephConfig{
			admSecret: map[string][]byte{
				"keyring": {0, 0, 0, 0},
			},
		}
		It("should return nil with secret", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
			}
			fakeClient := fake.NewFakeClient(secret)

			err := testConfig.UpdateKeyringToK8s(fakeClient, "test", "test")
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return err with no secret", func() {
			fakeClient := fake.NewFakeClient()
			err := testConfig.UpdateKeyringToK8s(fakeClient, "test", "test")
			Expect(err).To(HaveOccurred())
		})
	})
})
