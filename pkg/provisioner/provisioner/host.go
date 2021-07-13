package provisioner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/node"
)

func (p *Provisioner) applyHost(yamlWrapper wrapper.YamlInterface, execWrapper wrapper.ExecInterface, ioUtilWrapper wrapper.IoUtilInterface, cephConf, cephKeyring []byte) error {
	// Get host list to apply from CephCluster CR
	var err error
	var nodes []*node.Node
	var cephadmCurrentHostsBuf, hostAuthGetBuf bytes.Buffer

	cephHostsToApply := []*node.HostSpec{}
	cephHostNodesToApply := map[string]*node.Node{}

	nodes, err = p.getNodes()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		host := node.GetHostSpec()

		cephHostsToApply = append(cephHostsToApply, &host)
		hostName := host.GetHostName()

		cephHostNodesToApply[hostName] = node
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	// Get current ceph hosts
	cephName := p.getCephName()
	cephHostCheckCmd := []string{"orch", "host", "ls", "yaml"}

	fmt.Println("Executing: " + strings.Join(cephHostCheckCmd, ","))
	cephadmCurrentHostsBuf, err = util.RunCephCmd(wrapper.OsWrapper, execWrapper, ioUtilWrapper, cephConf, cephKeyring, cephName, cephHostCheckCmd...)
	if err != nil {
		fmt.Println("Error: " + cephadmCurrentHostsBuf.String())
		return err
	}

	fmt.Println("[applyHost] Existing hosts ---")
	fmt.Println(cephadmCurrentHostsBuf.String())

	// Extract host specs from ceph orch
	currentHosts := map[string]*node.HostSpec{}

	hostReader := bytes.NewReader(cephadmCurrentHostsBuf.Bytes())
	decoder := yamlWrapper.NewDecoder(hostReader)

	for {
		host := node.HostSpec{}
		if err = decoder.Decode(&host); err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		hostName := host.GetHostName()
		if err != nil {
			return err
		}
		currentHosts[hostName] = &host
	}

	pathConfigDir := p.getPathConfigDir()

	// generate public key
	hostAuthGetCmd := []string{"cephadm", "get-pub-key"}

	fmt.Println("Executing: " + strings.Join(hostAuthGetCmd, ","))
	hostAuthGetBuf, err = util.RunCephCmd(wrapper.OsWrapper, execWrapper, ioUtilWrapper, cephConf, cephKeyring, cephName, hostAuthGetCmd...)
	if err != nil {
		fmt.Println("Error: " + hostAuthGetBuf.String())
		return err
	}

	// Compare hosts in CR to Ceph and apply all changes
	for _, hostToApply := range cephHostsToApply {
		hostNameToApply := hostToApply.GetHostName()

		if _, exist := currentHosts[hostNameToApply]; exist {
			fmt.Println("Host EXIST!!!" + hostNameToApply)
			continue
		} else {
			// Write hostspec to yml
			hostFileName := fmt.Sprintf("%s%s.yml", pathConfigDir, hostNameToApply)
			fmt.Println("writing file to ", hostFileName)
			err = hostToApply.MakeYmlFile(yamlWrapper, ioUtilWrapper, hostFileName)
			if err != nil {
				return err
			}

			// Create generated key file
			pathHostPub := pathConfigDir + hostNameToApply + ".pub"
			err = ioUtilWrapper.WriteFile(pathHostPub, hostAuthGetBuf.Bytes(), 0644)
			if err != nil {
				return err
			}

			// Copy generated key
			var hostAuthApplyOutBuf, hostAuthApplyErrBuf, hostApplyBuf bytes.Buffer

			nodeID := cephHostNodesToApply[hostNameToApply].GetUserID()
			nodeIP := hostToApply.GetAddr()
			nodePw := cephHostNodesToApply[hostNameToApply].GetUserPw()

			const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"
			sshPassCmd := fmt.Sprintf("sshpass -f <(printf '%%s\\n' %s)", nodePw)
			hostAuthApplyCmd := fmt.Sprintf("%s ssh-copy-id %s -f -i %s %s@%s", sshPassCmd, sshKeyCheckOpt, pathHostPub, nodeID, nodeIP)

			fmt.Println("Executing: " + hostAuthApplyCmd)
			err = execWrapper.CommandExecute(ctx, &hostAuthApplyOutBuf, &hostAuthApplyErrBuf, "bash", "-c", hostAuthApplyCmd)
			if err != nil {
				fmt.Println("Error: " + hostAuthApplyErrBuf.String())
				return err
			}

			// Apply host
			hostApplyCmd := []string{"orch", "apply", "-i", hostFileName}

			fmt.Println("Executing: " + strings.Join(hostApplyCmd, ","))
			hostApplyBuf, err = util.RunCephCmd(wrapper.OsWrapper, execWrapper, ioUtilWrapper, cephConf, cephKeyring, cephName, hostAuthGetCmd...)

			if err != nil {
				fmt.Println("Error: " + hostApplyBuf.String())
				return err
			}

			fmt.Println(hostApplyBuf.String())
		}
	}

	// Check the result on ceph cluster hosts
	fmt.Println("Executing: " + strings.Join(cephHostCheckCmd, ","))
	cephadmCurrentHostsBuf, err = util.RunCephCmd(wrapper.OsWrapper, execWrapper, ioUtilWrapper, cephConf, cephKeyring, cephName, cephHostCheckCmd...)
	if err != nil {
		fmt.Println("Error: " + cephadmCurrentHostsBuf.String())
		return err
	}

	return nil
}
