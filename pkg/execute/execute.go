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
)

// Execution controlls command execution and stdout and stderr.
type Execution struct {
	cmd     *exec.Cmd
	stdout  *io.PipeReader
	stderr  *io.PipeReader
	closer  func()
	environ []string
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

// Stdout returns [io.Reader] for reading stdout.
func (e *Execution) Stdout() io.Reader { return e.stdout }

// Stderr returns [io.Reader] for reading stderr.
func (e *Execution) Stderr() io.Reader { return e.stderr }

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
	e.closer()

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
