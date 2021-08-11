package config

import (
	"context"
	"fmt"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CephConfig is struct for ceph config
type CephConfig struct {
	crConf    map[string]string
	admConf   map[string]string
	admSecret map[string][]byte
}

// SetCrConf sets crConf to value
func (conf *CephConfig) SetCrConf(m map[string]string) error {
	conf.crConf = m
	return nil
}

// SetAdmConf sets admConf to value
func (conf *CephConfig) SetAdmConf(m map[string]string) error {
	conf.admConf = m
	return nil
}

// SetAdmSecret sets admSecret to value
func (conf *CephConfig) SetAdmSecret(m map[string][]byte) error {
	conf.admSecret = m
	return nil
}

// GetCrConf gets value of crConf
func (conf *CephConfig) GetCrConf() map[string]string {
	return conf.crConf
}

// GetAdmConf gets value of admConf
func (conf *CephConfig) GetAdmConf() map[string]string {
	return conf.admConf
}

// GetAdmSecret gets value of admSecret
func (conf *CephConfig) GetAdmSecret() map[string][]byte {
	return conf.admSecret
}

// NewConfigFromCephCr reads ceph config from cephclusterspec and creates CephConfig
func NewConfigFromCephCr(cephCr hypersdsv1alpha1.CephClusterSpec) (*CephConfig, error) {
	conf := CephConfig{}
	crConf := make(map[string]string)

	for key, value := range cephCr.Config {
		crConf[key] = value
	}

	err := conf.SetCrConf(crConf)
	return &conf, err
}

// ConfigFromAdm reads ceph config(ceph access info) from ceph.conf
func (conf *CephConfig) ConfigFromAdm(ioUtil wrapper.IoUtilInterface, cephconf string) error {
	admConf := make(map[string]string)

	dat, err := ioUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}
	admConf["conf"] = string(dat)
	return conf.SetAdmConf(admConf)
}

// SecretFromAdm reads ceph keyring(ceph access info) from ceph.client.admin.keyring
func (conf *CephConfig) SecretFromAdm(ioUtil wrapper.IoUtilInterface, cephsecret string) error {
	admSecret := make(map[string][]byte)

	dat, err := ioUtil.ReadFile(cephsecret)
	if err != nil {
		return err
	}

	admSecret["keyring"] = dat
	return conf.SetAdmSecret(admSecret)
}

// MakeIniFile makes ceph.conf file using crConf
func (conf *CephConfig) MakeIniFile(ioUtil wrapper.IoUtilInterface, fileName string) error {
	ini := "[global]\n"
	crConf := conf.GetCrConf()
	for key, value := range crConf {
		s1 := fmt.Sprintf("\t%s = %s\n", key, value)
		ini = fmt.Sprintf("%s%s", ini, s1)
	}

	buf := []byte(ini)
	err := ioUtil.WriteFile(fileName, buf, 0644)

	return err
}

// UpdateConfToK8s updates ceph config(ceph access info) to k8s configmap
func (conf *CephConfig) UpdateConfToK8s(clientSet client.Client, cephNamespace, cephName string) error {
	configMap := &corev1.ConfigMap{}
	if err := clientSet.Get(context.TODO(), types.NamespacedName{Namespace: cephNamespace, Name: cephName + util.K8sConfigMapSuffix}, configMap); err != nil {
		return err
	}

	admConf := conf.GetAdmConf()
	configMap.Data = admConf
	err := clientSet.Update(context.TODO(), configMap)
	return err
}

// UpdateKeyringToK8s updates ceph keyring(ceph access info) to k8s secret
func (conf *CephConfig) UpdateKeyringToK8s(clientSet client.Client, cephNamespace, cephName string) error {
	secret := &corev1.Secret{}
	if err := clientSet.Get(context.TODO(), types.NamespacedName{Namespace: cephNamespace, Name: cephName + util.K8sSecretSuffix}, secret); err != nil {
		return err
	}

	admSecret := conf.GetAdmSecret()
	secret.Data = admSecret
	err := clientSet.Update(context.TODO(), secret)
	return err
}
