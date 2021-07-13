package util

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
)

const (
	// CephExecCmdTimeout defines max delay time for exec
	CephExecCmdTimeout = 1 * time.Minute
	// CephCmdTimeout defines max delay time for ceph
	CephCmdTimeout = "50"
)

// RunCephCmd extracts ceph access info from k8s configmap, secret and executes command to ceph
func RunCephCmd(os wrapper.OsInterface, exec wrapper.ExecInterface, ioUtil wrapper.IoUtilInterface,
	cephConf, cephKeyring []byte, cephName string, cmdQuery ...string) (bytes.Buffer, error) {
	var resultStdout, resultStderr bytes.Buffer

	pathConfigDir := PathGlobalConfigDir + cephName + "/"
	err := os.MkdirAll(pathConfigDir, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] CephName Driectory Create Error")
		return resultStdout, err
	}

	pathConf := pathConfigDir + CephConfName
	pathKeyring := pathConfigDir + CephKeyringName

	err = ioUtil.WriteFile(pathConf, cephConf, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] Ceph Conf Write Fail")
		return resultStdout, err
	}

	err = ioUtil.WriteFile(pathKeyring, cephKeyring, 0644)
	if err != nil {
		fmt.Println("[RunCephCmd] Ceph Keyring Write Fail")
		return resultStdout, err
	}

	cmdQuery = append(cmdQuery, "-c", pathConf, "--keyring", pathKeyring, "--connect-timeout", CephCmdTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), CephExecCmdTimeout)
	defer cancel()

	err = exec.CommandExecute(ctx, &resultStdout, &resultStderr, "ceph", cmdQuery...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}
