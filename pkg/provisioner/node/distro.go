package node

import (
	"errors"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"strconv"
	"strings"
)

// Distro contains linux distribution info of each node
type Distro string

const (
	// Ubuntu is a Linux distribution based on Debian
	Ubuntu Distro = "ubuntu"
	// Centos is a Linux distribution that functionally compatible with its upstream source, RHEL
	Centos Distro = "centos"
)

// Packager contains linux package manager info of each distro
type Packager string

const (
	// Yum is a primary package management tool for linux distro using rpm package manager
	Yum Packager = "yum "
	// Apt is a package management tool in Debian and Ubuntu
	Apt Packager = "apt-get "
)

// GetDistro fetch linux distribution info from the target node
func GetDistro(n *Node) error {
	distroStr, versionStr, err := fetchDistro(n)
	if err != nil {
		return err
	}
	packager, err := strToPackager(distroStr)
	if err != nil {
		return err
	}
	distro, err := strToDistro(distroStr)
	if err != nil {
		return err
	}
	version, err := strToVersion(versionStr)
	if err != nil {
		return err
	}

	err = n.SetOs(distro, packager, version)
	if err != nil {
		return err
	}

	return nil
}

func fetchDistro(n *Node) (distro, version string, err error) {
	cmd := "cat /etc/os-release | grep ID"
	output, err := n.RunSSHCmd(wrapper.SSHWrapper, cmd)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(output.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		str := strings.Split(line, "=")
		key := str[0]
		val := str[1]
		if val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if strings.EqualFold(key, "ID") {
			distro = val
		} else if strings.EqualFold(key, "VERSION_ID") {
			version = val
		}
	}

	return distro, version, nil
}

func strToPackager(distroName string) (Packager, error) {
	switch {
	case contains(getAptDistros(), distroName):
		return Apt, nil
	case contains(getYumDistros(), distroName):
		return Yum, nil
	default:
		return "", errors.New("OS distro is not supported: " + distroName)
	}
}

func strToDistro(distroName string) (Distro, error) {
	switch distroName {
	case string(Ubuntu):
		return Ubuntu, nil
	case string(Centos):
		return Centos, nil
	default:
		return "", errors.New("OS distro is not supported: " + distroName)
	}
}

func strToVersion(version string) (float64, error) {
	return strconv.ParseFloat(version, 64)
}

func getAptDistros() []string {
	return []string{string(Ubuntu)}
}

func getYumDistros() []string {
	return []string{string(Centos)}
}

func contains(arr []string, target string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}
