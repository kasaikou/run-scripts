package execute

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/usecases/execute"
)

// Execution controlls command execution and stdout and stderr.
type Execution struct {
	cmd          *exec.Cmd
	stdoutWriter *io.PipeWriter
	stdoutReader *io.PipeReader
	stderrWriter *io.PipeWriter
	stderrReader *io.PipeReader
	environ      []string
}

type executionBuilder func(execution models.Execution) (*Execution, error)

func (b executionBuilder) BuildExecution(execution models.Execution) (execute.Executor, error) {
	return b(execution)
}

// BuildExecution creates a instance from models.Execution.
//
// error will be models.CommandNotFoundError, models.EnvironmentNotFoundError or nil.
func BuildExecution(execution models.Execution) (*Execution, error) {
	path, args, err := createCmd(execution)
	if err != nil {
		return nil, err
	}

	if execution.Environments == nil {
		execution.Environments = []string{}
	}

	cmd := exec.Command(path, args...)
	cmd.Dir = execution.WorkingDir
	cmd.Env = execution.Environments

	return &Execution{
		cmd: cmd,
	}, nil
}

var ExecutionBuilder executionBuilder = BuildExecution

// Stdout returns [io.Reader] for reading stdout.
func (e *Execution) Stdout() io.Reader {
	if e.stdoutReader == nil {
		e.stdoutReader, e.stdoutWriter = io.Pipe()
		e.cmd.Stdout = e.stdoutWriter
	}

	return e.stdoutReader
}

// Stderr returns [io.Reader] for reading stderr.
func (e *Execution) Stderr() io.Reader {
	if e.stderrReader == nil {
		e.stderrReader, e.stderrWriter = io.Pipe()
		e.cmd.Stderr = e.stderrWriter
		e.stderrReader = e.stderrReader
	}

	return e.stderrReader
}

func (e *Execution) closePipe() {
	if e.stdoutWriter != nil {
		e.stdoutWriter.Close()
	}

	if e.stderrWriter != nil {
		e.stderrWriter.Close()
	}
}

func (e *Execution) Close() {
	if e.stdoutReader != nil {
		e.stdoutReader.Close()
	}

	if e.stderrReader != nil {
		e.stderrReader.Close()
	}
}

// Execute executes the configurated command.
//
// This function canonly be called once.
// It will panic if called more than once.
func (e *Execution) Execute(ctx context.Context) (exitCode int, err error) {

	if e.cmd.Process != nil {
		panic("execution has been already started")
	}

	tempEnvPath := filepath.Join(e.cmd.Dir, fmt.Sprintf("environ-%s", uuid.Must(uuid.NewV7()).String()))
	file, err := os.OpenFile(tempEnvPath, os.O_CREATE, 0600)
	if err != nil {
		return -1, fmt.Errorf("cannot open file: %w", err)
	}
	defer os.Remove(tempEnvPath)
	file.Close()

	e.cmd.Env = append(e.cmd.Env, "MARKFLOW_EXPORT="+tempEnvPath)

	err = e.cmd.Run()
	e.closePipe()

	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
		default:
			return -1, fmt.Errorf("execute command error: %w", err)
		}
	}

	file, err = os.Open(tempEnvPath)
	if err != nil {
		return -1, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return -1, fmt.Errorf("cannot read file: %w", err)
	}

	// TODO: Need to verify if it is a valid environment variables.
	environ := strings.Split(string(b), "\n")
	e.environ = make([]string, 0, len(environ))
	for _, env := range environ {
		if env != "" {
			e.environ = append(e.environ, env)
		}
	}

	return e.cmd.ProcessState.ExitCode(), nil
}

// Signal calls signal for running process.
//
// For windows, return models.NotImplementedForWindows.
// And If the process has been already exited, it will returns models.ErrExitedProcess.
func (e *Execution) Signal(sig os.Signal) error {
	if e.cmd.Process == nil {
		panic("execution has not been started yet")
	} else if runtime.GOOS == "windows" {
		return models.NotImplementedForWindows
	} else if e.cmd.ProcessState.Exited() {
		return models.ErrExitedProcess
	}

	return e.cmd.Process.Signal(sig)
}

// Kill kills the running process.
// If the process has been already exited, it will returns models.ErrExitedProcess.
func (e *Execution) Kill() error {
	if e.cmd.Process == nil {
		panic("execution has not been started yet")
	} else if e.cmd.ProcessState.Exited() {
		return models.ErrExitedProcess
	}

	return e.cmd.Process.Kill()
}

// Environ returns exported environment variables.
func (e *Execution) Environ() []string {
	return e.environ
}
