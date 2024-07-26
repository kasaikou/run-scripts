package execute_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kasaikou/markflow/pkg/execute"
	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/structured_log/trace"
	usecase_execute "github.com/kasaikou/markflow/pkg/usecases/execute"
	"github.com/kasaikou/markflow/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestMiddleIntegration_ExecuteExecution(t *testing.T) {

	wd := util.Must(os.Getwd())

	project := models.NewProject()
	executions := []models.Execution{
		models.Execution{
			ID:         uuid.Must(uuid.NewV7()),
			Name:       util.Must(models.ValidateExecutionName("target execution")),
			Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
			Script:     `echo "This is target execution"`,
			WorkingDir: wd,
			PrevExecutions: []models.ReferenceExecution{
				models.NewReferenceExecution(project, "parent execution 1"),
				models.NewReferenceExecution(project, "parent execution 2"),
			},
		},
		models.Execution{
			ID:         uuid.Must(uuid.NewV7()),
			Name:       util.Must(models.ValidateExecutionName("parent execution 1")),
			Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
			Script:     `echo "This is parent execution 1"`,
			WorkingDir: wd,
			PrevExecutions: []models.ReferenceExecution{
				models.NewReferenceExecution(project, "root execution"),
			},
		},
		models.Execution{
			ID:         uuid.Must(uuid.NewV7()),
			Name:       util.Must(models.ValidateExecutionName("parent execution 2")),
			Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
			Script:     `echo "This is parent execution 2"`,
			WorkingDir: wd,
			PrevExecutions: []models.ReferenceExecution{
				models.NewReferenceExecution(project, "root execution"),
			},
		},
		models.Execution{
			ID:         uuid.Must(uuid.NewV7()),
			Name:       util.Must(models.ValidateExecutionName("root execution")),
			Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
			Script:     `echo "This is root execution"`,
			WorkingDir: wd,
		},
	}

	for _, execution := range executions {
		require.NoError(t, project.AddExecution(execution))
	}

	request := usecase_execute.ExecuteExecutionRequest{
		Project:           project,
		CommandExecutions: []*models.Execution{&executions[0]},
		NumWorker:         2,
	}
	response := usecase_execute.ExecuteExecutionResponse{}

	ctx, _, close := trace.NewTracerAndRootSpan(context.Background(), "TestMiddleIntegration_ExecuteExecution")
	defer close()
	executor := usecase_execute.NewExecuteUseCase(slog.New(util.DiscardLogger), execute.ExecutionBuilder)
	executor.ExecuteExecution(ctx, request, &response)
}

func BenchmarkMiddleIntegration_ExecuteExecution(t *testing.B) {

	wd := util.Must(os.Getwd())

	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			project := models.NewProject()
			executions := []models.Execution{
				models.Execution{
					ID:         uuid.Must(uuid.NewV7()),
					Name:       util.Must(models.ValidateExecutionName("target execution")),
					Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
					Script:     `echo "This is target execution"`,
					WorkingDir: wd,
					PrevExecutions: []models.ReferenceExecution{
						models.NewReferenceExecution(project, "parent execution 1"),
						models.NewReferenceExecution(project, "parent execution 2"),
					},
				},
				models.Execution{
					ID:         uuid.Must(uuid.NewV7()),
					Name:       util.Must(models.ValidateExecutionName("parent execution 1")),
					Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
					Script:     `echo "This is parent execution 1"`,
					WorkingDir: wd,
					PrevExecutions: []models.ReferenceExecution{
						models.NewReferenceExecution(project, "root execution"),
					},
				},
				models.Execution{
					ID:         uuid.Must(uuid.NewV7()),
					Name:       util.Must(models.ValidateExecutionName("parent execution 2")),
					Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
					Script:     `echo "This is parent execution 2"`,
					WorkingDir: wd,
					PrevExecutions: []models.ReferenceExecution{
						models.NewReferenceExecution(project, "root execution"),
					},
				},
				models.Execution{
					ID:         uuid.Must(uuid.NewV7()),
					Name:       util.Must(models.ValidateExecutionName("root execution")),
					Lang:       util.Must(models.ValidateExecutionLanguage("sh")),
					Script:     `echo "This is root execution"`,
					WorkingDir: wd,
				},
			}

			for _, execution := range executions {
				require.NoError(t, project.AddExecution(execution))
			}

			request := usecase_execute.ExecuteExecutionRequest{
				Project:           project,
				CommandExecutions: []*models.Execution{&executions[0]},
				NumWorker:         2,
			}
			response := usecase_execute.ExecuteExecutionResponse{}

			ctx, _, close := trace.NewTracerAndRootSpan(context.Background(), "TestMiddleIntegration_ExecuteExecution")
			defer close()
			executor := usecase_execute.NewExecuteUseCase(slog.New(util.DiscardLogger), execute.ExecutionBuilder)
			executor.ExecuteExecution(ctx, request, &response)
		}
	})
}
