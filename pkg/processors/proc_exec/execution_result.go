package proc_exec

import (
	"sync"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
)

type ExecutionResult struct {
	Task   models.Task
	Status models.Status
}

var managerResultMapsMapping = [...]models.Status{
	models.StatusWaiting,
	models.StatusRunning,
	models.StatusSuccess,
	models.StatusFailed,
	models.StatusNotPlan,
}

type ExecutionResultsManager struct {
	lock     sync.RWMutex
	taskMaps []map[uuid.UUID]models.Task
}

func NewExecutionResultsManager() *ExecutionResultsManager {

	m := &ExecutionResultsManager{
		taskMaps: make([]map[uuid.UUID]models.Task, len(managerResultMapsMapping)),
	}

	for i := range len(m.taskMaps) {
		m.taskMaps[i] = make(map[uuid.UUID]models.Task)
	}

	return m
}

func (p *ExecutionResultsManager) statusIndex(status models.Status) int {
	for i := range managerResultMapsMapping {
		if managerResultMapsMapping[i] == status {
			return i
		}
	}

	panic("invalid status")
}

func (p *ExecutionResultsManager) Register(task models.Task, initStatus models.Status) {

	p.lock.Lock()
	func() {
		p.taskMaps[p.statusIndex(initStatus)][task.ID()] = task
	}()

	p.lock.Unlock()
}

func (p *ExecutionResultsManager) UpdateStatus(task models.Task, updateStatus models.Status) {

	p.lock.Lock()
	func() {
		for i := range p.taskMaps {
			task, exist := p.taskMaps[i][task.ID()]

			if exist {
				if updateStatus == managerResultMapsMapping[i] {
					return
				}

				delete(p.taskMaps[i], task.ID())
				p.taskMaps[p.statusIndex(updateStatus)][task.ID()] = task
			}
		}
	}()

	p.lock.Unlock()
}

func (p *ExecutionResultsManager) watchStatus(tasks ...models.Task) models.Status {

	notPlanTasks := p.taskMaps[p.statusIndex(models.StatusNotPlan)]
	if len(notPlanTasks) > 0 {
		for i := range tasks {
			if _, exist := notPlanTasks[tasks[i].ID()]; exist {
				return models.StatusNotPlan
			}
		}
	}

	failedTasks := p.taskMaps[p.statusIndex(models.StatusFailed)]
	if len(failedTasks) > 0 {
		for i := range tasks {
			if _, exist := failedTasks[tasks[i].ID()]; exist {
				return models.StatusFailed
			}
		}
	}

	allSuccess := true
	successTasks := p.taskMaps[p.statusIndex(models.StatusSuccess)]
	if len(successTasks) > 0 {
		for i := range tasks {
			if _, exist := successTasks[tasks[i].ID()]; !exist {
				allSuccess = false
			}
		}
	} else {
		allSuccess = false
	}

	if allSuccess {
		return models.StatusSuccess
	}

	allWaiting := true
	waitingTasks := p.taskMaps[p.statusIndex(models.StatusWaiting)]
	if len(waitingTasks) > 0 {
		for i := range tasks {
			if _, exist := waitingTasks[tasks[i].ID()]; !exist {
				allWaiting = false
			}
		}
	} else {
		allWaiting = false
	}

	if allWaiting {
		return models.StatusWaiting
	}

	return models.StatusRunning
}

func (p *ExecutionResultsManager) WatchStatus(tasks ...models.Task) models.Status {

	if len(tasks) == 0 {
		panic("executions argument is empty")
	}

	p.lock.RLock()
	result := p.watchStatus(tasks...)
	p.lock.RUnlock()
	return result
}
