package execution

import (
	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/models"
)

func ResolveAllDependencies(executions ...*models.Execution) []*models.Execution {

	depsMap := map[uuid.UUID]*models.Execution{}
	resolveAllDependencies(depsMap, executions)

	result := make([]*models.Execution, 0, len(depsMap))
	for _, execution := range depsMap {
		result = append(result, execution)
	}

	return result
}

func resolveAllDependencies(depsMap map[uuid.UUID]*models.Execution, executions []*models.Execution) {

	if executions == nil {
		return
	}

	for i := range executions {
		if _, exist := depsMap[executions[i].ID]; !exist {
			depsMap[executions[i].ID] = executions[i]

			deps := make([]*models.Execution, 0, len(executions[i].PrevExecutions))
			for j := range executions[i].PrevExecutions {
				deps = append(deps, executions[i].PrevExecutions[j].Execution())
			}

			resolveAllDependencies(depsMap, deps)
		}
	}
}
