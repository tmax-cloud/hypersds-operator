package node

import (
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"golang.org/x/crypto/ssh"

	"bytes"
	"context"
	"fmt"
)

type Node struct {
	userId   string
	userPw   string
	hostSpec HostSpec
}

func (n *Node) SetUserId(userId string) error {
	n.userId = userId
	return nil
}

func (n *Node) SetUserPw(userPw string) error {
	n.userPw = userPw
	return nil
}

func (n *Node) SetHostSpec(hostSpec HostSpec) error {
	n.hostSpec = hostSpec
	return nil
}

func (n *Node) GetUserId() string {
	return n.userId
}

func (n *Node) GetUserPw() string {
	return n.userPw
}

func (n *Node) GetHostSpec() HostSpec {
	return n.hostSpec
}

// TODO: replace sshpass command to go ssh pkg
func (n *Node) RunSshCmd(sshWrapper wrapper.SshInterface, cmdQuery string) (bytes.Buffer, error) {

	userPw := n.GetUserPw()
	userId := n.GetUserId()
	nodeHostSpec := n.GetHostSpec()
	ipAddr := nodeHostSpec.GetAddr()

	sshConfig := &ssh.ClientConfig{
		User: userId,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPw),
		},
		Timeout:         SshCmdTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint    , todo ssh key verification
	}

	var resultStdout, resultStderr bytes.Buffer
	err := sshWrapper.Run(ipAddr, cmdQuery, &resultStdout, &resultStderr, sshConfig)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}

/* Executing command
 * (DESTINATION) sshpass -f <(printf '%s\n' userPw) scp -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null srcFile userId@ipAddr:/destFile
 * (SOURCE) sshpass -f <(printf '%s\n' userPw) scp -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null userId@ipAddr:/srcFile destFile
 */
// TODO: replace sshpass command to go ssh pkg
func (n *Node) RunScpCmd(exec wrapper.ExecInterface, srcFile, destFile string, role Role) (bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), SshCmdTimeout)
	defer cancel()

	userPw := n.GetUserPw()
	userId := n.GetUserId()
	nodeHostSpec := n.GetHostSpec()
	ipAddr := nodeHostSpec.GetAddr()

	const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"

	var scpCmd string
	// provisioner sends srcFile to this node as destFile
	if role == DESTINATION {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[3]s %[4]s@%[5]s:/%[6]s",
			userPw, sshKeyCheckOpt, srcFile, userId, ipAddr, destFile)

		// this node sends srcFile to provisioner as destFile
	} else {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[4]s@%[5]s:%[3]s %[6]s",
			userPw, sshKeyCheckOpt, srcFile, userId, ipAddr, destFile)
	}

	parameters := []string{"-c"}
	parameters = append(parameters, scpCmd)

	var resultStdout, resultStderr bytes.Buffer
	err := exec.CommandExecute(&resultStdout, &resultStderr, ctx, "bash", parameters...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}

func NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]*Node, error) {
	var nodes []*Node

	for _, nodeInCephSpec := range cephSpec.Nodes {
		var n Node
		err := n.SetUserId(nodeInCephSpec.UserID)
		if err != nil {
			return nil, err
		}
		err = n.SetUserPw(nodeInCephSpec.Password)
		if err != nil {
			return nil, err
		}

		var hostSpec HostSpec
		err = hostSpec.SetServiceType()
		if err != nil {
			return nil, err
		}

		err = hostSpec.SetAddr(nodeInCephSpec.IP)
		if err != nil {
			return nil, err
		}

		err = hostSpec.SetHostName(nodeInCephSpec.HostName)
		if err != nil {
			return nil, err
		}

		err = n.SetHostSpec(hostSpec)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}
