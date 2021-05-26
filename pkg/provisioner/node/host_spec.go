package node

import (
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"

	"errors"
)

// HostSpec is struct to define ceph host spec of ceph spec yaml
type HostSpec struct {
	// ServiceType indicates service_type field
	ServiceType string `yaml:"service_type"` // use const var
	// Addr indicates addr field
	Addr string `yaml:"addr"`
	// HostName indicates hostname field
	HostName string `yaml:"hostname"`
	// Labels indicates labels field
	Labels []string `yaml:"labels,omitempty"`
	// Status indicates status field
	Status string `yaml:"status,omitempty"`
}

// TODO: Add validations (ex. IP check, ..)

// SetServiceType sets ServiceType to value
func (hs *HostSpec) SetServiceType() error {
	hs.ServiceType = HostSpecServiceType
	return nil
}

// SetAddr sets Addr to value
func (hs *HostSpec) SetAddr(addr string) error {
	hs.Addr = addr
	return nil
}

// SetHostName sets HostName to value
func (hs *HostSpec) SetHostName(hostName string) error {
	if hostName == "" {
		return errors.New("HostName must not be empty string")
	}

	hs.HostName = hostName
	return nil
}

// SetLabels sets Labels to value
func (hs *HostSpec) SetLabels(labels []string) error {
	hs.Labels = append([]string{}, labels...)
	return nil
}

// AddLabels adds value to Labels
func (hs *HostSpec) AddLabels(labels ...string) error {
	hs.Labels = append(hs.Labels, labels...)
	return nil
}

// SetStatus sets Status to value
func (hs *HostSpec) SetStatus(status string) error {
	hs.Status = status
	return nil
}

// GetServiceType gets HostSpecServiceTypee
func (hs *HostSpec) GetServiceType() string {
	return HostSpecServiceType
}

// GetHostName gets value of HostName
func (hs *HostSpec) GetHostName() string {
	return hs.HostName
}

// GetAddr gets value of Addr
func (hs *HostSpec) GetAddr() string {
	return hs.Addr
}

// GetLabels gets value of Labels
func (hs *HostSpec) GetLabels() []string {
	return hs.Labels
}

// GetStatus gets value of Status
func (hs *HostSpec) GetStatus() string {
	return hs.Status
}

// MakeYmlFile creates ceph host service file using HostSpec
func (hs *HostSpec) MakeYmlFile(yamlWrapper wrapper.YamlInterface, ioUtilWrapper wrapper.IoUtilInterface, fileName string) error {
	ymlBytes, err := yamlWrapper.Marshal(hs)
	if err != nil {
		return err
	}

	err = ioUtilWrapper.WriteFile(fileName, ymlBytes, 0644)

	return err
}
