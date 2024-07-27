package execute

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"sync"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/structured_log/trace"
	"github.com/kasaikou/markflow/pkg/util"
)

type ExecutionWorker struct {
	NumWorker        int
	Logger           *slog.Logger
	executeEventChan chan ExecuteEventRequest
	lock             sync.Mutex
	isAlreadyRouted  bool
	isClosed         bool
}

func NewExecutionWorker() *ExecutionWorker {
	return &ExecutionWorker{
		NumWorker:        runtime.NumCPU(),
		Logger:           slog.Default(),
		executeEventChan: make(chan ExecuteEventRequest),
	}
}

func (e *ExecutionWorker) Route(ctx context.Context) {

	func() {
		e.lock.Lock()
		defer e.lock.Unlock()
		if e.isAlreadyRouted {
			panic("execution worker is already routed")
		}
		e.isAlreadyRouted = true
	}()

	defer close(e.executeEventChan)
	defer func() { e.isClosed = true }()
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for i := range e.NumWorker {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.routeWorker(ctx, i)
		}()
	}
}

func (e *ExecutionWorker) SendExecuteRequest(req ExecuteEventRequest) {

	func() {
		e.lock.Lock()
		defer e.lock.Unlock()
		if !e.isAlreadyRouted {
			panic("execution worker is not routed")
		} else if e.isClosed {
			panic("already worker is closed")
		}

		e.executeEventChan <- req
	}()

}

type ExecuteEventRequest struct {
	Execution *models.Execution
	Context   context.Context
	OnStart   func()
	OnSucceed func()
	OnFail    func()
}

func (e *ExecutionWorker) routeWorker(ctx context.Context, idx int) {

	recv := e.executeEventChan
	for {
		select {
		case <-ctx.Done():
			return

		case req := <-recv:
			e.execute(&req, idx)
		}
	}
}

func (e *ExecutionWorker) execute(req *ExecuteEventRequest, idx int) {

	execution := req.Execution
	ctx, closeSpan := trace.WithSpan(req.Context, strings.Join([]string{"execute/worker/executions", execution.ID.String()}, "/"), slog.Int("worker_idx", idx))
	defer closeSpan()

	builtExec, err := BuildExecution(*execution)
	if err != nil {
		e.Logger.ErrorContext(ctx, "Error in Building Execution.", slog.Any("error", err))
		util.FuncOrDefault(req.OnFail, nil)
	}

	e.Logger.InfoContext(ctx, "Start Executing Command.")
	util.FuncOrDefault(req.OnStart, nil)
	exitCode, err := builtExec.Execute(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "Error in Executing Command.", slog.Any("error", err))
		util.FuncOrDefault(req.OnFail, nil)
	} else if exitCode != 0 {
		e.Logger.ErrorContext(ctx, "Exit code is not unexpected.", slog.Int("want", 0), slog.Int("have", exitCode))
		util.FuncOrDefault(req.OnFail, nil)
	}

	e.Logger.InfoContext(ctx, "Succeeded for Executing Command.")
	util.FuncOrDefault(req.OnSucceed, nil)
}
