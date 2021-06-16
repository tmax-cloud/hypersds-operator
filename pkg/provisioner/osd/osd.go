package osd

import (
	"bytes"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	"github.com/tmax-cloud/hypersds-operator/pkg/common/wrapper"
	"github.com/tmax-cloud/hypersds-operator/pkg/provisioner/service"
)

// Spec is struct to define ceph osd spec of ceph spec yaml
type Spec struct {
	// DataDevices indicates data_devices field
	DataDevices Device `yaml:"data_devices,omitempty"`
	// DbDevices indicates db_devices field
	DbDevices Device `yaml:"db_devices,omitempty"`
	// WalDevices indicates wal_devices field
	WalDevices Device `yaml:"wal_devices,omitempty"`
	// JournalDevices indicates journal_devices field
	JournalDevices Device `yaml:"journal_devices,omitempty"`
	// DataDirectories indicates data_directories field
	DataDirectories []string `yaml:"data_directories,omitempty"`
	// OsdsPerDevice indicates osds_per_device field
	OsdsPerDevice int `yaml:"osds_per_device,omitempty"`
	// Objectstore indicates objectstore field
	Objectstore string `yaml:"objectstore,omitempty"`
	// Encrypted indicates encrypted field
	Encrypted bool `yaml:"encrypted,omitempty"`
	// FilterLogic indicates filter_logic field
	FilterLogic string `yaml:"filter_logic,omitempty"`
	/*
	    				db_slots=None,  # type: Optional[int]
	                    wal_slots=None,  # type: Optional[int]
	                    osd_id_claims=None,  # type: Optional[Dict[str, List[str]]]
	                    block_db_size=None,  # type: Union[int, str, None]
	                    block_wal_size=None,  # type: Union[int, str, None]
	                    journal_size=None,  # type: Union[int, str, None]
	                    service_type=None,  # type: Optional[str]
	                    unmanaged=False,  # type: bool
	   				preview_only=False,  # type: bool
	*/
}

// Osd is struct to define ceph osd service of ceph spec yaml
type Osd struct {
	// Service defines ceph service spec of ceph spec yaml
	Service service.Service `yaml:",inline"`
	// Spec defines ceph osd spec of ceph spec yaml
	Spec Spec `yaml:"spec,omitempty"`
}

// SetService sets Service to value
func (o *Osd) SetService(s *service.Service) error {
	o.Service = *s
	return nil
}

// SetDataDevices sets DataDevices in Spec to value
func (o *Osd) SetDataDevices(dataDevices *Device) error {
	o.Spec.DataDevices = *dataDevices
	return nil
}

// GetService gets value of Service
func (o *Osd) GetService() service.Service {
	return o.Service
}

// GetDataDevices gets value of DataDevices in Spec
func (o *Osd) GetDataDevices() Device {
	return o.Spec.DataDevices
}

// CompareDataDevices compares between Osd and targetOsd and returns addDeviceList,removeDeviceList
func (o *Osd) CompareDataDevices(targetOsd *Osd) (addDeviceList, removeDeviceList []string, err error) {
	// o: orch osd, targetOsd: cephCr osd
	dataDevices := o.GetDataDevices()

	devicePaths := dataDevices.getPaths()

	targetDataDevices := targetOsd.GetDataDevices()

	targetDevicePaths := targetDataDevices.getPaths()

	deviceMap := map[string]bool{}

	for _, device := range devicePaths {
		deviceMap[device] = false
	}
	for _, device := range targetDevicePaths {
		_, exists := deviceMap[device]
		if exists {
			deviceMap[device] = true
		} else {
			addDeviceList = append(addDeviceList, device)
		}
	}

	for device, value := range deviceMap {
		if !value {
			removeDeviceList = append(removeDeviceList, device)
		}
	}
	return addDeviceList, removeDeviceList, nil
}

// MakeYmlFile creates ceph osd service file using Osd
func (o *Osd) MakeYmlFile(yaml wrapper.YamlInterface, ioUtilWrapper wrapper.IoUtilInterface, fileName string) error {
	osdYaml, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	err = ioUtilWrapper.WriteFile(fileName, osdYaml, 0644)
	return err
}

// NewOsdsFromCephCr reads osd information from cephclusterspec and creates Osd list
func NewOsdsFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]*Osd, error) {
	var osds []*Osd

	for _, Spec := range cephSpec.Osd {
		var hosts []string
		var osd Osd
		var dataDevices Device
		var s service.Service
		var placement service.Placement

		// set Placement, Service
		hosts = append(hosts, Spec.HostName)
		err := placement.SetHosts(hosts)
		if err != nil {
			return nil, err
		}
		err = s.SetPlacement(placement)
		if err != nil {
			return nil, err
		}
		err = s.SetServiceType("osd")
		if err != nil {
			return nil, err
		}
		err = s.SetServiceID("osd_" + Spec.HostName)
		if err != nil {
			return nil, err
		}
		// set device
		err = dataDevices.setPaths(Spec.Devices)
		if err != nil {
			return nil, err
		}
		err = osd.SetDataDevices(&dataDevices)
		if err != nil {
			return nil, err
		}
		err = osd.SetService(&s)
		if err != nil {
			return nil, err
		}
		osds = append(osds, &osd)
	}

	return osds, nil
}

// NewOsdsFromCephOrch reads osd information from ceph orch and creates Osd list
func NewOsdsFromCephOrch(yaml wrapper.YamlInterface, rawOsdsFromOrch []byte) ([]*Osd, error) {
	var osds []*Osd
	readerOrch := bytes.NewReader(rawOsdsFromOrch)
	dec := yaml.NewDecoder(readerOrch)
	for {
		var osd Osd
		err := dec.Decode(&osd)
		if err != nil {
			break
		}
		osds = append(osds, &osd)
	}
	return osds, nil
}
