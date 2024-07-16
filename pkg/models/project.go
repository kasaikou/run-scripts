package models

import (
	"encoding/json"
	"fmt"
)

// Project contains execution and pipelines.
type Project struct {
	executions     []Execution
	executionNames map[string]*Execution
	pipelines      []Pipeline
	pipelineNames  map[string]*Pipeline
}

// NewProject creates a Project instance.
func NewProject() *Project {
	return &Project{
		executionNames: make(map[string]*Execution),
		pipelineNames:  make(map[string]*Pipeline),
	}
}

// AddExecution adds a Execution.
func (p *Project) AddExecution(ex Execution) error {
	if _, duplicated := p.executionNames[ex.Name.String()]; duplicated {
		return NewModelValidateError(fmt.Errorf("execution name '%s' is duplicated", ex.Name.String()))
	}
	for _, alias := range ex.Aliases {
		if _, duplicated := p.executionNames[alias]; duplicated {
			return NewModelValidateError(fmt.Errorf("the alias '%s' of execution '%s' is duplicated", alias, ex.Name.String()))
		}
	}

	p.executions = append(p.executions, ex)
	p.executionNames[ex.Name.String()] = &p.executions[len(p.executions)-1]
	for _, alias := range ex.Aliases {
		p.executionNames[alias] = &p.executions[len(p.executions)-1]
	}

	return nil
}

// AddPipeline adds a Pipeline.
func (p *Project) AddPipeline(pl Pipeline) error {
	if _, duplicated := p.pipelineNames[pl.Name.String()]; duplicated {
		return NewModelValidateError(fmt.Errorf("pipeline name '%s' is duplicated", pl.Name.String()))
	}
	for _, alias := range pl.Aliases {
		if _, duplicated := p.pipelineNames[alias]; duplicated {
			return NewModelValidateError(fmt.Errorf("the alias '%s' of pipeline '%s' is duplicated", alias, pl.Name.String()))
		}
	}

	p.pipelines = append(p.pipelines, pl)
	p.pipelineNames[pl.Name.String()] = &p.pipelines[len(p.pipelines)-1]
	for _, alias := range pl.Aliases {
		p.pipelineNames[alias] = &p.pipelines[len(p.pipelines)-1]
	}

	return nil
}

// Executions provides iterator for Execution.
func (p *Project) Executions() func(yield func(int, *Execution) bool) {
	return func(yield func(int, *Execution) bool) {
		for i := 0; i < len(p.executions); i++ {
			if !yield(i, &p.executions[i]) {
				return
			}
		}
	}
}

// Pipelines provides iterator for registered Pipeline.
func (p *Project) Pipelines() func(yield func(int, *Pipeline) bool) {
	return func(yield func(int, *Pipeline) bool) {
		for i := 0; i < len(p.pipelines); i++ {
			if !yield(i, &p.pipelines[i]) {
				return
			}
		}
	}
}

// GetExecution gets the Execution instance associated with the name.
func (p *Project) GetExecution(name string) *Execution {
	if execution, exist := p.executionNames[name]; exist {
		return execution
	}
	return nil
}

// GetPipeline gets the Pipeline instance associated with the name.
func (p *Project) GetPipeline(name string) *Pipeline {
	if pipeline, exist := p.pipelineNames[name]; exist {
		return pipeline
	}
	return nil
}

// ProjectJSONContent is a structure used to expand Project in JSON format.
type ProjectJSONContent struct {
	Executions []ExecutionJSONContent `json:"executions"`
	Pipelines  []PipelineJSONContent  `json:"pipelines"`
}

// ToJSONContent converts the instance to ProjectJSONContent.
func (p *Project) ToJSONContent() ProjectJSONContent {

	executions := make([]ExecutionJSONContent, 0, len(p.executions))
	for _, execution := range p.executions {
		executions = append(executions, execution.ToJSONContent())
	}

	pipelines := make([]PipelineJSONContent, 0, len(p.pipelines))
	for _, pipeline := range p.pipelines {
		pipelines = append(pipelines, pipeline.ToJSONContent())
	}

	return ProjectJSONContent{
		Executions: executions,
		Pipelines:  pipelines,
	}
}

// MarshalJSON represents the instance in JSON format.
func (p *Project) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToJSONContent())
}
