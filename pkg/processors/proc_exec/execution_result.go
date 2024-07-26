package proc_exec

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
)

type ExecutionResult struct {
	Execution *models.Execution
	Status    models.ExecutionStatus
}

type ExecutionResultsManager struct {
	lock      sync.RWMutex
	resultMap map[uuid.UUID]ExecutionResult
}

func NewExecutionResultsManager() *ExecutionResultsManager {
	return &ExecutionResultsManager{
		resultMap: make(map[uuid.UUID]ExecutionResult),
	}
}

func (p *ExecutionResultsManager) Register(execution *models.Execution, initStatus models.ExecutionStatus) {

	p.lock.Lock()
	func() {
		p.resultMap[execution.ID] = ExecutionResult{
			Execution: execution,
			Status:    initStatus,
		}
	}()

	p.lock.Unlock()
}

func (p *ExecutionResultsManager) UpdateStatus(execution *models.Execution, updateStatus models.ExecutionStatus) {

	p.lock.Lock()
	func() {
		result, exist := p.resultMap[execution.ID]
		if !exist {
			return
		}

		result.Status = updateStatus
		p.resultMap[execution.ID] = result
	}()

	p.lock.Unlock()
}

func (p *ExecutionResultsManager) watchStatus(executions ...*models.Execution) models.ExecutionStatus {

	result := p.resultMap[executions[0].ID]
	if !result.Status.IsValid() {
		panic(fmt.Sprintf("status of '%s' execution is undefined", result.Execution.Name.String()))
	}

	switch result.Status {
	case models.ExecFailed, models.ExecNotPlan:
		return result.Status
	}

	status := result.Status

	for i := 1; i < len(executions); i++ {
		result := p.resultMap[executions[i].ID]
		if !result.Status.IsValid() {
			panic(fmt.Sprintf("status of '%s' execution is undefined", result.Execution.Name.String()))
		}

		switch result.Status {
		case models.ExecFailed, models.ExecNotPlan:
			return result.Status

		case models.ExecRunning:
			switch status {
			case models.ExecSuccess, models.ExecWaiting:
				status = models.ExecRunning
			}

		case models.ExecWaiting:
			switch status {
			case models.ExecSuccess, models.ExecRunning:
				status = models.ExecRunning
			}
		}
	}

	return status
}

func (p *ExecutionResultsManager) WatchStatus(executions ...*models.Execution) models.ExecutionStatus {

	if len(executions) == 0 {
		panic("executions argument is empty")
	}

	p.lock.RLock()
	result := p.watchStatus(executions...)
	p.lock.RUnlock()
	return result
}

func (p *ExecutionResultsManager) PreparedExecution() *models.Execution {

	p.lock.RLock()
	execution := func() *models.Execution {
		deps := make([]*models.Execution, 0, len(p.resultMap))

		for _, result := range p.resultMap {
			if result.Status == models.ExecWaiting {
				if len(result.Execution.PrevExecutions) == 0 {
					return result.Execution
				}

				deps := deps[:0]
				for i := range result.Execution.PrevExecutions {
					deps = append(deps, result.Execution.PrevExecutions[i].Execution())
				}

				status := p.watchStatus(deps...)
				if status == models.ExecSuccess {
					return result.Execution
				} else if status == models.ExecFailed {
					return nil
				}

			} else if result.Status == models.ExecFailed {
				return nil
			}
		}

		return nil
	}()
	p.lock.RUnlock()

	return execution
}

// func (p *ExecutionResultsManager) Contains(status models.ExecutionStatus, executions ...*models.Execution) bool {

// 	p.lock.RLock()
// 	result := func() bool {
// 		for i := range executions {
// 			result := p.resultMap[executions[i].ID]
// 			if result.Status == status {
// 				return true
// 			} else if !result.Status.IsValid() {
// 				panic(fmt.Sprintf("status of '%s' execution is undefined", result.Execution.Name.String()))
// 			}
// 		}

// 		return false
// 	}()
// 	p.lock.RUnlock()
// 	return result
// }

// func (p *ExecutionResultsManager) Every(status models.ExecutionStatus, executions ...*models.Execution) bool {

// 	p.lock.RLock()
// 	result := func() bool {
// 		for i := range executions {
// 			result := p.resultMap[executions[i].ID]
// 			if result.Status != status {
// 				return false
// 			}
// 		}

// 		return true
// 	}()
// 	p.lock.RUnlock()
// 	return result
// }
