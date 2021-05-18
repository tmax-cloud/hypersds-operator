package osd

// https://github.com/ceph/ceph/blob/master/src/python-common/ceph/deployment/drive_group.py

type Device struct {
	Paths      []string `yaml:"paths,omitempty"`
	Model      string   `yaml:"model,omitempty"`
	Size       string   `yaml:"size,omitempty"`
	Rotational bool     `yaml:"rotational,omitempty"`
	Limit      int      `yaml:"limit,omitempty"`
	Vendor     string   `yaml:"vendor,omitempty"`
	All        bool     `yaml:"all,omitempty"`
}

func (d *Device) SetPaths(paths []string) error {
	d.Paths = paths
	return nil
}
func (d *Device) SetModel(model string) error {
	d.Model = model
	return nil
}
func (d *Device) SetSize(size string) error {
	d.Size = size
	return nil
}
func (d *Device) SetRotational(rotational bool) error {
	d.Rotational = rotational
	return nil
}
func (d *Device) SetLimit(limit int) error {
	d.Limit = limit
	return nil
}
func (d *Device) SetVendor(vendor string) error {
	d.Vendor = vendor
	return nil
}
func (d *Device) SetAll(all bool) error {
	d.All = all
	return nil
}

func (d *Device) GetPaths() []string {
	return d.Paths
}
func (d *Device) GetModel() string {
	return d.Model
}
func (d *Device) GetSize() string {
	return d.Size
}
func (d *Device) GetRotational() bool {
	return d.Rotational
}
func (d *Device) GetLimit() int {
	return d.Limit
}
func (d *Device) GetVendor() string {
	return d.Vendor
}
func (d *Device) GetAll() bool {
	return d.All
}
