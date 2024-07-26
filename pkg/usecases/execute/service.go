package execute

import (
	"context"
	"io"
	"os"

	"github.com/kasaikou/markflow/pkg/models"
)

// ExecutionBuilder is the interface for building from models.Execution to Executor.
type ExecutionBuilder interface {
	BuildExecution(execution models.Execution) (Executor, error)
}

// Executor is the interface for executing process.
type Executor interface {
	Execute(ctx context.Context) (exitCode int, err error)
	Stdout() io.Reader
	Stderr() io.Reader
	Signal(sig os.Signal) error
	Kill() error
}
