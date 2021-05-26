package wrapper

import (
	"bytes"
	"context"
	"os/exec"
)

// ExecInterface is wrapper interface of exec package
type ExecInterface interface {
	// CommandExecute is wrapper function for CommandContext and Run function of exec package
	CommandExecute(ctx context.Context, resultStdout, resultStderr *bytes.Buffer, name string, arg ...string) error
}

// execStruct is struct to be used instead of exec package
type execStruct struct {
}

// CommandExecute executes CommandContext and Run function of exec package
func (e *execStruct) CommandExecute(ctx context.Context, resultStdout, resultStderr *bytes.Buffer, name string, arg ...string) error {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdout = resultStdout
	cmd.Stderr = resultStderr

	return cmd.Run()
}

// ExecWrapper is to be used instead of exec package
var ExecWrapper ExecInterface

func init() {
	ExecWrapper = &execStruct{}
}
