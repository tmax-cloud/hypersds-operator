package node

import (
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"golang.org/x/crypto/ssh"

	"bytes"
	"context"
	"fmt"
)

// Node contains node info to deploy ceph
type Node struct {
	userID   string
	userPw   string
	hostSpec HostSpec
}

// SetUserID sets userID to value
func (n *Node) SetUserID(userID string) error {
	n.userID = userID
	return nil
}

// SetUserPw sets userPw to value
func (n *Node) SetUserPw(userPw string) error {
	n.userPw = userPw
	return nil
}

// SetHostSpec sets hostSpec to value
func (n *Node) SetHostSpec(hostSpec *HostSpec) error {
	n.hostSpec = *hostSpec
	return nil
}

// GetUserID gets value of userID
func (n *Node) GetUserID() string {
	return n.userID
}

// GetUserPw gets value of userPw
func (n *Node) GetUserPw() string {
	return n.userPw
}

// GetHostSpec gets value of hostSpec
func (n *Node) GetHostSpec() HostSpec {
	return n.hostSpec
}

// RunSSHCmd executes the command on the node using ssh package
func (n *Node) RunSSHCmd(sshWrapper wrapper.SSHInterface, cmdQuery string) (bytes.Buffer, error) {
	userPw := n.GetUserPw()
	userID := n.GetUserID()
	nodeHostSpec := n.GetHostSpec()
	ipAddr := nodeHostSpec.GetAddr()

	sshConfig := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPw),
		},
		Timeout:         SSHCmdTimeout,
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
 * (DESTINATION) sshpass -f <(printf '%s\n' userPw) scp -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null srcFile userID@ipAddr:/destFile
 * (SOURCE) sshpass -f <(printf '%s\n' userPw) scp -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null userID@ipAddr:/srcFile destFile
 */
// TODO: replace sshpass command to go ssh pkg

// RunScpCmd executes file copy between source node and destination node using scp
func (n *Node) RunScpCmd(exec wrapper.ExecInterface, srcFile, destFile string, role Role) (bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), SSHCmdTimeout)
	defer cancel()

	userPw := n.GetUserPw()
	userID := n.GetUserID()
	nodeHostSpec := n.GetHostSpec()
	ipAddr := nodeHostSpec.GetAddr()

	const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"

	var scpCmd string
	// provisioner sends srcFile to this node as destFile
	if role == DESTINATION {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[3]s %[4]s@%[5]s:%[6]s",
			userPw, sshKeyCheckOpt, srcFile, userID, ipAddr, destFile)

		// this node sends srcFile to provisioner as destFile
	} else {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[4]s@%[5]s:%[3]s %[6]s",
			userPw, sshKeyCheckOpt, srcFile, userID, ipAddr, destFile)
	}

	parameters := []string{"-c"}
	parameters = append(parameters, scpCmd)

	var resultStdout, resultStderr bytes.Buffer
	err := exec.CommandExecute(ctx, &resultStdout, &resultStderr, "bash", parameters...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}

// NewNodesFromCephCr reads node information from cephclusterspec and creates Node list
func NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]*Node, error) {
	var nodes []*Node

	for _, nodeInCephSpec := range cephSpec.Nodes {
		var n Node
		err := n.SetUserID(nodeInCephSpec.UserID)
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

		err = n.SetHostSpec(&hostSpec)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}
