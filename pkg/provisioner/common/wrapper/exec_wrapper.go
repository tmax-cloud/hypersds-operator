package wrapper

import (
	"bytes"
	"context"
	"os/exec"
)

// exec package의 function 대신 호출될 exec interface
type ExecInterface interface {
	CommandExecute(resultStdout, resultStderr *bytes.Buffer, ctx context.Context, name string, arg ...string) error
}

// exec interface method를 가진 struct
// 실제 코드에서는 해당 struct의 method가 호출되고,
// unit test 시에는 이 것과 비슷한 구조의 mock struct를 구현하여(gomock이 자동으로 생성해줌) mock struct의 method가 호출된다.
type execStruct struct {
}

// 실제 코드에서 exec interface를 통해 수행되는 method
// exec package function 수행을 감싸는 method
func (e *execStruct) CommandExecute(resultStdout, resultStderr *bytes.Buffer, ctx context.Context, name string, arg ...string) error {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdout = resultStdout
	cmd.Stderr = resultStderr

	return cmd.Run()
}

// package 변수로 ExecWrapper를 정의하여 main.go나 hypersds-provisioner.go 같은 곳에서 sshCommand function 호출시 ExecWrapper를 넘겨주도록 함
// output, err := util.RunSSHCmd(util.ExecWrapper, hostName, hostAddr, testCommand...)
var ExecWrapper ExecInterface

func init() {
	//ExecWrapper의 ExecStruct를 할당하여 실제 코드 상에서 최종적으로 ExecStruct.commandExecute를 호출되게 함
	ExecWrapper = &execStruct{}
}
