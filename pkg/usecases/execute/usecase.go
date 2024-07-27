package execute

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/processors/execution"
	"github.com/kasaikou/markflow/pkg/processors/proc_exec"
	"github.com/kasaikou/markflow/pkg/structured_log/trace"
	"github.com/kasaikou/markflow/pkg/util"
)

type ExecuteUseCase struct {
	logger      *slog.Logger
	execBuilder ExecutionBuilder
}

type ExecuteExecutionRequest struct {
	Project           *models.Project
	CommandExecutions []*models.Execution
	NumWorker         int
}

type ExecuteExecutionResponse struct {
	Results *proc_exec.ExecutionResultsManager
}

func NewExecuteUseCase(logger *slog.Logger, execBuilder ExecutionBuilder) *ExecuteUseCase {
	return &ExecuteUseCase{
		logger:      logger,
		execBuilder: execBuilder,
	}
}

type executeExecutionWorkerRequest struct {
	manager   *proc_exec.ExecutionResultsManager
	execution *models.Execution
	onStart   func()
	onSucceed func()
	onFail    func()
}

func (u *ExecuteUseCase) ExecuteExecution(ctx context.Context, req ExecuteExecutionRequest, res *ExecuteExecutionResponse) {

	executions := execution.ResolveAllDependencies(req.CommandExecutions...)
	results := proc_exec.NewExecutionResultsManager()
	for i := range executions {
		results.Register(executions[i], models.StatusWaiting)
	}
	res.Results = results

	requestChan := make(chan executeExecutionWorkerRequest)
	defer close(requestChan)
	startedChan := make(chan struct{})
	defer close(startedChan)
	stateChangeChan := make(chan struct{}, len(executions)*2)
	defer close(stateChangeChan)

	wgControllers := sync.WaitGroup{}
	defer wgControllers.Wait()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wgControllers.Add(1)
	go func() {
		defer wgControllers.Done()
		for {
			for {
				execution := results.PreparedExecution()
				if execution != nil {
					requestChan <- executeExecutionWorkerRequest{
						manager:   results,
						execution: execution,
						onStart:   func() { startedChan <- struct{}{} },
						onSucceed: func() { stateChangeChan <- struct{}{} },
						onFail: func() {
							cancel()
							stateChangeChan <- struct{}{}
						},
					}
					<-startedChan
					continue
				}

				switch status := results.WatchStatus(executions...); status {
				case models.StatusSuccess, models.StatusFailed:
					cancel()
					return
				}

				break
			}

			<-stateChangeChan
			for len(stateChangeChan) > 0 {
				<-stateChangeChan
			}
		}
	}()

	wgWorkers := sync.WaitGroup{}
	defer wgWorkers.Wait()
	for i := range min(req.NumWorker, len(executions)) {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			u.executeExecutionWorker(ctx, i, requestChan)
		}()
	}

	wgWorkers.Wait()
}

func (u *ExecuteUseCase) executeExecutionWorker(ctx context.Context, idx int, recv <-chan executeExecutionWorkerRequest) {
	workerLogger := u.logger.With(slog.Int("worker_idx", idx))

	for req := range util.RouteChanContext(ctx, recv) {
		ctx, closeSpan := trace.WithSpan(ctx, strings.Join([]string{"execute/worker/executions", req.execution.ID.String()}, "/"))
		defer closeSpan()

		logger := workerLogger

		execution, err := u.execBuilder.BuildExecution(*req.execution)
		if err != nil {
			req.manager.UpdateStatus(req.execution, models.StatusFailed)
			logger.ErrorContext(ctx, "Error in Building Execution.", slog.Any("error", err))
			req.onFail()
		}

		req.manager.UpdateStatus(req.execution, models.StatusRunning)
		logger.InfoContext(ctx, "Start Executing Command.")
		req.onStart()
		exitCode, err := execution.Execute(ctx)

		if err != nil {
			req.manager.UpdateStatus(req.execution, models.StatusFailed)
			logger.ErrorContext(ctx, "Error in Executing Command.", slog.Any("error", err))
			req.onFail()
		} else if exitCode != 0 {
			req.manager.UpdateStatus(req.execution, models.StatusFailed)
			logger.ErrorContext(ctx, "Exit code is not 0.", slog.Int("exit_code", exitCode))
			req.onFail()
		}

		req.manager.UpdateStatus(req.execution, models.StatusSuccess)
		logger.InfoContext(ctx, "Suceeded for Executing Command.")
		req.onSucceed()
	}
}
