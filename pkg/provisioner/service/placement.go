package service

// Placement is struct for placement field of ceph service spec yaml
type Placement struct {
	// Label indicates label field
	Label string `yaml:"label,omitempty"`
	// Hosts indicates host field
	Hosts []string `yaml:"hosts,omitempty"`
	// Count indicates count field
	Count int `yaml:"count,omitempty"`
	// HostPattern indicates host_pattern field
	HostPattern string `yaml:"host_pattern,omitempty"`
}

// SetLabel sets Label to value
func (p *Placement) SetLabel(label string) error {
	p.Label = label
	return nil
}

// SetHosts sets Hosts to value
func (p *Placement) SetHosts(hosts []string) error {
	p.Hosts = hosts
	return nil
}

// SetCount sets Count to value
func (p *Placement) SetCount(count int) error {
	p.Count = count
	return nil
}

// SetHostPattern sets HostPattern to value
func (p *Placement) SetHostPattern(hostPattern string) error {
	p.HostPattern = hostPattern
	return nil
}

// GetLabel gets value of Label
func (p Placement) GetLabel() string {
	return p.Label
}

// GetHosts gets value of Hosts
func (p Placement) GetHosts() []string {
	return p.Hosts
}

// GetCount gets value of Count
func (p Placement) GetCount() int {
	return p.Count
}

// GetHostPattern gets value of HostPattern
func (p Placement) GetHostPattern() string {
	return p.HostPattern
}
