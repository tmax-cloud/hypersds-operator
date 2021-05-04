package wrapper

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

const (
	Port    = "22"
	Network = "tcp"
)

type SshInterface interface {
	Run(addr, command string, resultStdout, resultStderr *bytes.Buffer, config *ssh.ClientConfig) error
}

type sshStruct struct {
}

func (s *sshStruct) Run(addr, command string, resultStdout, resultStderr *bytes.Buffer, config *ssh.ClientConfig) error {
	target := addr + ":" + Port
	client, err := ssh.Dial(Network, target, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = resultStdout
	session.Stderr = resultStderr
	return session.Run(command)
}

var SshWrapper SshInterface

func init() {
	SshWrapper = &sshStruct{}
}
