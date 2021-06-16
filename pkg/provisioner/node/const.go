package node

import (
	"time"
)

const (
	// HostSpecServiceType indicates host service of ceph spec yaml
	HostSpecServiceType = "host"
	// SSHCmdTimeout  defines max delay time for SSH connection
	SSHCmdTimeout = 20 * time.Minute
)

// Role indicates target of SSH connection
type Role int

const (
	// SOURCE indicates source of SSH connection
	SOURCE Role = iota
	// DESTINATION indicates destination of SSH connection
	DESTINATION
)
