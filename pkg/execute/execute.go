package execute

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/kasaikou/markflow/pkg/models"
)

type Execution struct {
	cmd    *exec.Cmd
	stdout *io.PipeReader
	stderr *io.PipeReader
	closer func()
}

func BuildExecution(execution models.Execution) (*Execution, error) {
	path, args, err := createCmd(execution)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, args...)
	cmd.Dir = execution.WorkingDir
	cmd.Env = execution.Environ

	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	return &Execution{
		cmd:    cmd,
		stdout: stdoutReader,
		stderr: stderrReader,
		closer: func() {
			stdoutWriter.Close()
			stderrWriter.Close()
		},
	}, nil
}

func (e *Execution) Stdout() io.Reader { return e.stdout }
func (e *Execution) Stderr() io.Reader { return e.stderr }

func (e *Execution) Execute(ctx context.Context) (exitCode int, err error) {

	if e.cmd.Process != nil {
		panic("execution has been already started")
	}

	err = e.cmd.Run()
	e.closer()

	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
		default:
			return -1, fmt.Errorf("execute command error: %w", err)
		}
	}

	return e.cmd.ProcessState.ExitCode(), nil
}
