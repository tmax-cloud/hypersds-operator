package node

import (
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	"errors"
)

// variables are required to be importable so that yaml wrapper marshal/unmarshal them
type HostSpec struct {
	ServiceType string   `yaml:"service_type"` // use const var
	Addr        string   `yaml:"addr"`
	HostName    string   `yaml:"hostname"`
	Labels      []string `yaml:"labels,omitempty"`
	Status      string   `yaml:"status,omitempty"`
}

// TODO: Add validations (ex. IP check, ..)

func (hs *HostSpec) SetServiceType() error {
	hs.ServiceType = HostSpecServiceType
	return nil
}

func (hs *HostSpec) SetAddr(addr string) error {
	hs.Addr = addr
	return nil
}

func (hs *HostSpec) SetHostName(hostName string) error {
	if hostName == "" {
		return errors.New("HostName must not be empty string")
	}

	hs.HostName = hostName
	return nil
}

func (hs *HostSpec) SetLabels(labels []string) error {
	hs.Labels = append([]string{}, labels...)
	return nil
}

func (hs *HostSpec) AddLabels(labels ...string) error {
	hs.Labels = append(hs.Labels, labels...)
	return nil
}

func (hs *HostSpec) SetStatus(status string) error {
	hs.Status = status
	return nil
}

func (hs *HostSpec) GetServiceType() string {
	return HostSpecServiceType
}

func (hs *HostSpec) GetHostName() string {
	return hs.HostName
}

func (hs *HostSpec) GetAddr() string {
	return hs.Addr
}

func (hs *HostSpec) GetLabels() []string {
	return hs.Labels
}

func (hs *HostSpec) GetStatus() string {
	return hs.Status
}

func (hs *HostSpec) MakeYmlFile(yamlWrapper wrapper.YamlInterface, ioUtilWrapper wrapper.IoUtilInterface, fileName string) error {
	ymlBytes, err := yamlWrapper.Marshal(hs)
	if err != nil {
		return err
	}

	err = ioUtilWrapper.WriteFile(fileName, ymlBytes, 0644)

	return err
}
