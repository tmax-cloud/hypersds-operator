package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	corev1 "k8s.io/api/core/v1"
)

const (
	CephExecCmdTimeout = 1 * time.Minute
	CephCmdTimeout     = "50"
)

func RunCephCmd(os wrapper.OsInterface, exec wrapper.ExecInterface, ioUtil wrapper.IoUtilInterface, cephConf *corev1.ConfigMap, cephKeyring *corev1.Secret, cephName string, cmdQuery ...string) (bytes.Buffer, error) {
	var resultStdout, resultStderr bytes.Buffer

	pathConfigDir := PathGlobalConfigDir + cephName + "/"
	err := os.MkdirAll(pathConfigDir, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] CephName Driectory Create Error")
		return resultStdout, err
	}

	pathConf := pathConfigDir + CephConfName
	pathKeyring := pathConfigDir + CephKeyringName

	confData, ok := cephConf.Data["conf"]
	if !ok {
		fmt.Println("[RunCephCmd] Ceph Conf Data isn't exist")
		return resultStdout, errors.New("[RunCephCmd] Ceph Conf Data isn't exist")
	}
	confDataBuf := []byte(confData)
	err = ioUtil.WriteFile(pathConf, confDataBuf, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] Ceph Conf Write Fail")
		return resultStdout, err
	}

	keyringDataBuf, ok := cephKeyring.Data["keyring"]
	if !ok {
		fmt.Println("[RunCephCmd] Ceph Keyring Data isn't exist")
		return resultStdout, errors.New("[RunCephCmd] Ceph keyring Data isn't exist")
	}
	err = ioUtil.WriteFile(pathKeyring, keyringDataBuf, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] Ceph Keyring Write Fail")
		return resultStdout, err
	}

	cmdQuery = append(cmdQuery, "-c", pathConf, "--keyring", pathKeyring, "--connect-timeout", CephCmdTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), CephExecCmdTimeout)
	defer cancel()

	err = exec.CommandExecute(&resultStdout, &resultStderr, ctx, "ceph", cmdQuery...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}
