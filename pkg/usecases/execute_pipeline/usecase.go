package execute_pipeline

import (
	"log/slog"

	"github.com/kasaikou/markflow/pkg/models"
	"github.com/kasaikou/markflow/pkg/usecases/execute"
)

type ExecutePipelineUseCase struct {
	logger      *slog.Logger
	execBuilder execute.ExecutionBuilder
}

type ExecutePipeline struct {
	Project  models.Project
	Pipeline *models.Pipeline
}
