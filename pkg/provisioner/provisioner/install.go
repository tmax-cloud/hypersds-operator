package provisioner

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/node"
)

const cmdNotFound = "command not found"

func (p *Provisioner) updateCephClusterToOp() error {
	fmt.Println("\n----------------Start to update conf and keyring to operator---------------")

	cephConfig := p.getCephConfig()
	clientSet := p.getClientSet()
	cephNamespace := p.getCephNamespace()
	cephName := p.getCephName()
	pathConf := p.getPathConfigDir() + util.CephConfName
	pathKeyring := p.getPathConfigDir() + util.CephKeyringName

	err := cephConfig.ConfigFromAdm(wrapper.IoUtilWrapper, pathConf)
	if err != nil {
		return err
	}

	err = cephConfig.SecretFromAdm(wrapper.IoUtilWrapper, pathKeyring)
	if err != nil {
		return err
	}

	err = cephConfig.UpdateConfToK8s(clientSet, cephNamespace, cephName)
	if err != nil {
		return err
	}

	err = cephConfig.UpdateKeyringToK8s(clientSet, cephNamespace, cephName)

	return err
}

func isCephadmInstalled(deployNode *node.Node) (bool, error) {
	const checkCephadmInstalledCmd = "cephadm version"
	output, err := deployNode.RunSSHCmd(wrapper.SSHWrapper, checkCephadmInstalledCmd)
	if err != nil {
		if strings.Contains(output.String(), cmdNotFound) {
			return false, nil
		}
		fmt.Println("[Provisioner] cephadm installation check is failed")
		return false, err
	}
	return true, nil
}

func installCephadm(targetNode *node.Node) error {
	installed, err := isCephadmInstalled(targetNode)
	if err != nil {
		return err
	}
	if installed {
		return nil
	}
	if err := installCephadmPackage(targetNode); err != nil {
		return err
	}
	return nil
}

func installPackages(nodeList []*node.Node) error {
	if err := fetchOSInfo(nodeList); err != nil {
		return err
	}

	for _, n := range nodeList {
		installed, err := isDockerInstalled(n)
		if installed {
			continue
		}
		if err != nil {
			return err
		}
		err = installCommonPackage(n)
		if err != nil {
			return err
		}
		err = installDocker(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func fetchOSInfo(targetNodeList []*node.Node) error {
	for _, n := range targetNodeList {
		if err := node.GetDistro(n); err != nil {
			return err
		}
	}

	return nil
}

func processCmdOnNode(targetNode *node.Node, command string) error {
	output, err := targetNode.RunSSHCmd(wrapper.SSHWrapper, command)
	return processExecError(err, output)
}

func createDstDir(targetNode *node.Node, wholePath string) error {
	paths := strings.Split(wholePath, "/")
	paths = paths[:len(paths)-1]
	dirPath := ""

	for _, path := range paths {
		dirPath += path + "/"
	}
	dirPath = strings.TrimSuffix(dirPath, "/")
	cmd := "mkdir -p " + dirPath

	return processCmdOnNode(targetNode, cmd)
}

func copyFile(targetNode *node.Node, role node.Role, srcFile, destFile string) error {
	output, err := targetNode.RunScpCmd(wrapper.ExecWrapper, srcFile, destFile, role)
	return processExecError(err, output)
}

func processExecError(errExec error, output bytes.Buffer) error {
	if errExec != nil {
		if output.Bytes() != nil {
			_, err := output.WriteTo(os.Stderr)
			if err != nil {
				// TODO: combine errExec and err
				return err
			}
		}
		return errExec
	}

	_, err := output.WriteTo(os.Stdout)

	return err
}
