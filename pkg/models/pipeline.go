package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// PipelineName represents pipeline name.
type PipelineName struct{ name string }

// ValidatePipelineName checks name is a valid pipeline name and create a PipelineName instance.
func ValidatePipelineName(name string) (PipelineName, error) {
	return PipelineName{name: name}, nil
}

// String indicates a pipeline name.
func (p PipelineName) String() string { return p.name }

// PipelineStep contains configuration of pipeline's step.
type PipelineStep struct {
	Name       string
	Executions []ReferenceExecution
}

// PipelineStepJSONContent is a structure used to expand PipelineStep in JSON format.
type PipelineStepJSONContent struct {
	Name       string   `json:"name"`
	Executions []string `json:"executes"`
}

// ToJSONContent converts the instance to PipelineStepJSONContent.
func (p *PipelineStep) ToJSONContent() PipelineStepJSONContent {

	executions := make([]string, 0, len(p.Executions))
	for _, execution := range p.Executions {
		executions = append(executions, execution.String())
	}

	return PipelineStepJSONContent{
		Name:       p.Name,
		Executions: executions,
	}
}

// Pipeline contains configurations of execution pipeline.
type Pipeline struct {
	ID      uuid.UUID
	Name    PipelineName
	Aliases []string
	Steps   []PipelineStep
}

// NewPipeline creates a Pipeline instance.
func NewPipeline() Pipeline {
	return Pipeline{ID: uuid.Must(uuid.NewV7())}
}

// PipelineJSONContent is a structure used to expand Pipeline in JSON format.
type PipelineJSONContent struct {
	Name    string                    `json:"name"`
	Aliases []string                  `json:"aliases"`
	Steps   []PipelineStepJSONContent `json:"steps"`
}

// ToJSONContent converts the instance to PipelineJSONContent.
func (p Pipeline) ToJSONContent() PipelineJSONContent {

	steps := make([]PipelineStepJSONContent, 0, len(p.Steps))
	for _, step := range p.Steps {
		steps = append(steps, step.ToJSONContent())
	}

	return PipelineJSONContent{
		Name:    p.Name.String(),
		Aliases: p.Aliases,
		Steps:   steps,
	}
}

// MarshalJSON represents the instance in JSON format.
func (p *Pipeline) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToJSONContent())
}
