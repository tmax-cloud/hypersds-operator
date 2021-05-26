package wrapper

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

const (
	// Port is default port to be used of ssh communication
	port = "22"
	// Network is default protocol to be used of ssh communication
	network = "tcp"
)

// SSHInterface is wrapper interface of ssh package
type SSHInterface interface {
	// Run is wrapper function for ssh execution process of ssh package
	Run(addr, command string, resultStdout, resultStderr *bytes.Buffer, config *ssh.ClientConfig) error
}

// sshStruct is struct to be used instead of ssh package
type sshStruct struct {
}

func deferCheck(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Defer error:", err)
	}
}

// Run executes ssh execution process of ssh package
func (s *sshStruct) Run(addr, command string, resultStdout, resultStderr *bytes.Buffer, config *ssh.ClientConfig) error {
	target := addr + ":" + port
	client, err := ssh.Dial(network, target, config)
	if err != nil {
		return err
	}
	defer deferCheck(client.Close)

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer deferCheck(session.Close)
	session.Stdout = resultStdout
	session.Stderr = resultStderr
	return session.Run(command)
}

// SSHWrapper is to be used instead of ssh package
var SSHWrapper SSHInterface

func init() {
	SSHWrapper = &sshStruct{}
}
