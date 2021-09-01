package provisioner

import (
	"fmt"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/util"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/node"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/osd"
	"strings"
)

const (
	installCmd = "install -y "
	updateCmd  = "update -y"
)

func isDockerInstalled(targetNode *node.Node) (bool, error) {
	const checkDockerCmd = "docker -v"
	output, err := targetNode.RunSSHCmd(wrapper.SSHWrapper, checkDockerCmd)
	if err != nil {
		if strings.Contains(output.String(), cmdNotFound) {
			return false, nil
		}
		fmt.Println("[Provisioner] docker installation check is failed")
		return false, err
	}
	return true, nil
}

func installCommonPackage(n *node.Node) error {
	fmt.Println("\n----------------Start to install base package---------------")
	var err error
	packager := n.GetOs().Packager

	fmt.Println("[installBasePackage] executing packager update")
	if packager == node.Apt {
		err = processCmdOnNode(n, string(packager)+updateCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] install common packages")
	const commonPkgs = "chrony python3"
	err = processCmdOnNode(n, string(packager)+installCmd+commonPkgs)
	if err != nil {
		return err
	}

	fmt.Println("[installBasePackage] enable chronyd")
	const enableChronydCmd = "systemctl enable chronyd"
	err = processCmdOnNode(n, enableChronydCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installBasePackage] restart chronyd")
	const restartChronydCmd = "systemctl restart chronyd"
	err = processCmdOnNode(n, restartChronydCmd)
	if err != nil {
		return err
	}

	return nil
}

func installDocker(n *node.Node) error {
	fmt.Println("[installBasePackage] install docker dependencies")
	var err error
	packager := n.GetOs().Packager

	const dockerDepAptPkgs = "apt-transport-https ca-certificates curl gnupg lsb-release"
	const dockerDepYumPkgs = "yum-utils"
	pkgs := dockerDepAptPkgs
	if packager == node.Yum {
		pkgs = dockerDepYumPkgs
	}
	err = processCmdOnNode(n, string(packager)+installCmd+pkgs)
	if err != nil {
		return err
	}

	fmt.Println("[installBasePackage] executing curl docker ...")
	const dockerUbuntuGpgKeyCmd = "curl -s https://download.docker.com/linux/ubuntu/gpg | apt-key add - &>/dev/null"
	osName := n.GetOs().Distro
	if osName == node.Ubuntu {
		err = processCmdOnNode(n, dockerUbuntuGpgKeyCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing add-apt-repo docker ...")
	const dockerAptRepoCmd = `add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`
	const dockerYumRepoCmd = "yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo"
	addDockerRepoCmd := ""
	switch packager {
	case node.Apt:
		addDockerRepoCmd += dockerAptRepoCmd
	case node.Yum:
		addDockerRepoCmd += dockerYumRepoCmd
	}
	err = processCmdOnNode(n, addDockerRepoCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installBasePackage] executing install docker-ce")
	const dockerPkgs = "docker-ce"
	if packager == node.Apt {
		err = processCmdOnNode(n, string(packager)+updateCmd)
		if err != nil {
			return err
		}
	}
	err = processCmdOnNode(n, string(packager)+installCmd+dockerPkgs)
	if err != nil {
		return err
	}

	fmt.Println("[installBasePackage] executing sysctl docker")
	const restartDockerCmd = "systemctl restart docker"
	err = processCmdOnNode(n, restartDockerCmd)
	if err != nil {
		return err
	}

	return nil
}

func installCephadmPackage(targetNode *node.Node) error {
	fmt.Println("\n----------------Start to install cephadm---------------")

	fmt.Println("[installCephadm] executing curl cephadm")
	curlCephadmCmd := fmt.Sprintf("curl --silent --remote-name --location https://github.com/ceph/ceph/raw/v%s/src/cephadm/cephadm", cephVersion)
	err := processCmdOnNode(targetNode, curlCephadmCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing chmod")
	const chmodCmd = "chmod +x cephadm"
	err = processCmdOnNode(targetNode, chmodCmd)
	if err != nil {
		return err
	}

	// TODO: Specify release version
	fmt.Println("[installCephadm] executing cephadm add-repo")
	admAddRepoCmd := fmt.Sprintf("./cephadm add-repo --version %s", cephVersion)
	err = processCmdOnNode(targetNode, admAddRepoCmd)
	if err != nil {
		return err
	}

	if targetNode.GetOs().Distro == node.Ubuntu {
		// need for ubuntu 18.04 & ceph version 15.2.8
		fmt.Println("[installCephadm] executing curl cephadm gpg key")
		addCephadmRepoCmd := "curl https://download.ceph.com/keys/release.asc | " +
			"gpg --no-default-keyring --keyring /tmp/fix.gpg --import - && " +
			"gpg --no-default-keyring --keyring /tmp/fix.gpg --export > /etc/apt/trusted.gpg.d/ceph.release.gpg && rm /tmp/fix.gpg"
		err = processCmdOnNode(targetNode, addCephadmRepoCmd)
		if err != nil {
			return err
		}

		fmt.Println("[installCephadm] executing cephadm apt-get update")
		const aptUpdateCmd = "apt-get update"
		err = processCmdOnNode(targetNode, aptUpdateCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installCephadm] executing cephadm install")
	const admInstallCmd = "./cephadm install"
	err = processCmdOnNode(targetNode, admInstallCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing mkdir")
	const mkdirCmd = "mkdir -p /etc/ceph"
	err = processCmdOnNode(targetNode, mkdirCmd)

	return err
}

func bootstrapCephadm(targetNode *node.Node, pathConfFromCr string) error {
	fmt.Println("\n----------------Start to bootstrap ceph---------------")

	fmt.Println("[bootstrapCephadm] copying conf file")
	err := createDstDir(targetNode, pathConfFromCr)
	if err != nil {
		return err
	}
	err = copyFile(targetNode, node.DESTINATION, pathConfFromCr, pathConfFromCr)
	if err != nil {
		return err
	}

	deployNodeHostSpec := targetNode.GetHostSpec()

	monIP := deployNodeHostSpec.GetAddr()

	fmt.Println("[bootstrapCephadm] executing bootstrap")
	admBootstrapCmd := fmt.Sprintf("cephadm --image %s bootstrap --mon-ip %s --config %s",
		cephImageName, monIP, pathConfFromCr)
	err = processCmdOnNode(targetNode, admBootstrapCmd)
	if err != nil {
		return err
	}

	fmt.Println("[bootstrapCephadm] checking status")
	const admHealthCheckCmd = "cephadm shell -- ceph -s"
	err = processCmdOnNode(targetNode, admHealthCheckCmd)

	return err
}

func (p *Provisioner) applyOsd(cephConf, cephKeyring []byte) error {
	var err error

	fmt.Println("[applyOsd] get osds from CephOrch")

	cephName := p.getCephName()
	pathConfigDir := p.getPathConfigDir()

	cmd := []string{"orch", "ls", "--service_type", "osd", "--export", "--refresh"}
	output, err := util.RunCephCmd(wrapper.OsWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring, cephName, cmd...)
	if err != nil {
		return processExecError(err, output)
	}

	var osdsFromOrch []*osd.Osd
	if !strings.Contains(output.String(), "No services reported") {
		rawOsdsFromOrch := output.Bytes()

		osdsFromOrch, err = osd.NewOsdsFromCephOrch(wrapper.YamlWrapper, rawOsdsFromOrch)
		if err != nil {
			return err
		}
	}
	osdsFromCephCr, err := osd.NewOsdsFromCephCr(p.cephCluster)
	if err != nil {
		return err
	}

	var osdMap map[string]*osd.Osd
	var removeOsdMap map[string]bool

	osdMap = make(map[string]*osd.Osd)
	removeOsdMap = make(map[string]bool)

	for _, osdOrch := range osdsFromOrch {
		osdService := osdOrch.GetService()

		osdServiceID := osdService.GetServiceID()

		osdMap[osdServiceID] = osdOrch
		removeOsdMap[osdServiceID] = true
	}

	fmt.Println("[applyOsd] compare osds between CephCR and CephOrch")

	for _, osdCephCr := range osdsFromCephCr {
		osdService := osdCephCr.GetService()
		osdServiceID := osdService.GetServiceID()
		osdOrch, exist := osdMap[osdServiceID]

		var addDeviceList, removeDeviceList []string

		if exist {
			addDeviceList, removeDeviceList, err = osdOrch.CompareDataDevices(osdCephCr)
			if err != nil {
				return err
			}
			removeOsdMap[osdServiceID] = false
			// todo remove disk ....
			fmt.Printf("[applyOsd] osd service: %s, add: %+q, remove: %+q\n", osdServiceID, addDeviceList, removeDeviceList)
		}

		fmt.Println("[applyOsd] make osd yaml")

		osdFileName := pathConfigDir + osdServiceID + ".yaml"
		err = osdCephCr.MakeYmlFile(wrapper.YamlWrapper, wrapper.IoUtilWrapper, osdFileName)
		if err != nil {
			return err
		}

		fmt.Printf("[applyOsd] apply osd service: %s\n", osdServiceID)

		applyCmd := []string{"orch", "apply", "-i", osdFileName}
		output, err = util.RunCephCmd(wrapper.OsWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring, cephName, applyCmd...)
		if err != nil {
			return processExecError(err, output)
		}
	}
	for osdServiceID, value := range removeOsdMap {
		if !value {
			continue
		}
		osdServiceName := "osd." + osdServiceID

		fmt.Printf("[applyOsd] remove osd service: %s\n", osdServiceName)

		removeCmd := []string{"orch", "rm", osdServiceName}
		output, err = util.RunCephCmd(wrapper.OsWrapper, wrapper.ExecWrapper, wrapper.IoUtilWrapper, cephConf, cephKeyring, cephName, removeCmd...)
		if err != nil {
			return processExecError(err, output)
		}
	}
	return nil
}
