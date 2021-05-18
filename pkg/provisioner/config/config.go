package config

import (
	"context"
	"fmt"
	"strings"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (conf *CephConfig) SetCrConf(m map[string]string) error {
	conf.crConf = m
	return nil
}

func (conf *CephConfig) SetAdmConf(m map[string]string) error {
	conf.admConf = m
	return nil
}

func (conf *CephConfig) SetAdmSecret(m map[string][]byte) error {
	conf.admSecret = m
	return nil
}

func (conf *CephConfig) GetCrConf() map[string]string {
	return conf.crConf
}

func (conf *CephConfig) GetAdmConf() map[string]string {
	return conf.admConf
}

func (conf *CephConfig) GetAdmSecret() map[string][]byte {
	return conf.admSecret
}

type CephConfig struct {
	crConf    map[string]string
	admConf   map[string]string
	admSecret map[string][]byte
}

func NewConfigFromCephCr(cephCr hypersdsv1alpha1.CephClusterSpec) (*CephConfig, error) {
	conf := CephConfig{}
	crConf := make(map[string]string)

	for key, value := range cephCr.Config {
		crConf[key] = value
	}

	err := conf.SetCrConf(crConf)
	return &conf, err
}

func (conf *CephConfig) ConfigFromAdm(ioUtil wrapper.IoUtilInterface, cephconf string) error {
	admConf := make(map[string]string)

	dat, err := ioUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}
	lines := strings.Split(string(dat[:]), "\n")
	for _, s := range lines {
		kv := strings.Split(s, "=")
		if len(kv) > 1 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			admConf[key] = val
		}
	}
	admConf["conf"] = string(dat)
	return conf.SetAdmConf(admConf)
}

func (conf *CephConfig) SecretFromAdm(ioUtil wrapper.IoUtilInterface, cephsecret string) error {
	admSecret := make(map[string][]byte)

	dat, err := ioUtil.ReadFile(cephsecret)
	if err != nil {
		return err
	}

	admSecret["keyring"] = dat
	return conf.SetAdmSecret(admSecret)
}

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

//default:default SA가 configmap put할수 있는 rolebinding 있어야함!
func (conf *CephConfig) UpdateConfToK8s(clientSet client.Client, cephNamespace, cephName string) error {
	configMap := &corev1.ConfigMap{}
	if err := clientSet.Get(context.TODO(), types.NamespacedName{Namespace: cephNamespace, Name: cephName}, configMap); err != nil {
		return err
	}

	admConf := conf.GetAdmConf()
	configMap.Data = admConf
	err := clientSet.Update(context.TODO(), configMap)
	return err
}

func (conf *CephConfig) UpdateKeyringToK8s(clientSet client.Client, cephNamespace, cephName string) error {
	secret := &corev1.Secret{}
	if err := clientSet.Get(context.TODO(), types.NamespacedName{Namespace: cephNamespace, Name: cephName}, secret); err != nil {
		return err
	}

	admSecret := conf.GetAdmSecret()
	secret.Data = admSecret
	err := clientSet.Update(context.TODO(), secret)
	return err
}
