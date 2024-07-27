package execute

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/util"
)

func newExecuteTask(execution *models.Execution, worker *ExecutionWorker, taskReferer models.TaskReferer) *executeTask {
	return &executeTask{}
}

type executeTask struct {
	execution *models.Execution
	worker    *ExecutionWorker
}

func (t *executeTask) Kind() string { return "execution" }

func (t *executeTask) ID() uuid.UUID { return t.execution.ID }

func (t *executeTask) Start(ctx context.Context, req models.StartEventRequest) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	t.worker.SendExecuteRequest(ExecuteEventRequest{
		Execution: t.execution,
		Context:   ctx,
		OnStart: func() {
			defer wg.Done()
			util.FuncOrDefault(req.OnStart, nil)
		},
		OnSucceed: req.OnSucceed,
		OnFail:    req.OnFail,
	})
}

func (t *executeTask) Depends() []models.Task {
	t.execution.PrevExecutions
	return
}
