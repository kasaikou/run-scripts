package models

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// ExecutionName represents execution name.
type ExecutionName struct{ name string }

// ValidateExecutionName checks name is a valid execution name and create a ExecutionName instance.
func ValidateExecutionName(name string) (ExecutionName, error) {
	return ExecutionName{name: name}, nil
}

// String indicates execution name.
func (en ExecutionName) String() string { return en.name }

// ExecutionPath represents execution path.
//
// If set to an empty string, it means that this is the execution path set in PATH.
type ExecutionPath struct{ path string }

// ValidateExecutionPath checks path is a valid execution path and create ExecutionPath instance.
func ValidateExecutionPath(path string) (ExecutionPath, error) {
	return ExecutionPath{path: path}, nil
}

func (ep ExecutionPath) String() string { return ep.path }

// ExecutionLanguage represents execution language.
type ExecutionLanguage struct{ lang string }

// RegExecutionLanguage is the validation rule for execution language.
var RegExecutionLanguage = regexp.MustCompile(`^[a-z]+$`)

// ValidateExecutionLanguage checks lang is a validation execution language and create ExecutionLanguage instance.
func ValidateExecutionLanguage(lang string) (ExecutionLanguage, error) {

	if !RegExecutionLanguage.MatchString(lang) {
		return ExecutionLanguage{}, NewModelValidateError(fmt.Errorf("'%s' is not satisfied regular expression ('%s')", lang, RegExecutionLanguage.String()))
	}

	return ExecutionLanguage{lang: lang}, nil
}

// String indicates execution language.
func (el ExecutionLanguage) String() string { return el.lang }

// Execution contains command execution configuration.
type Execution struct {
	ID             uuid.UUID
	Name           ExecutionName
	Descriptions   []string
	Aliases        []string
	Path           ExecutionPath
	PrevExecutions []ReferenceExecution
	Lang           ExecutionLanguage
	Script         string
	Environments   []string
	WorkingDir     string
	AdditionalArgs []string
}

// NewExecution creates a Execution instance.
func NewExecution() Execution {
	return Execution{ID: uuid.Must(uuid.NewV7())}
}

// ExecutionJSONContent is a strucuture used to expand Execution in JSON format.
type ExecutionJSONContent struct {
	Name               string   `json:"name"`
	Path               string   `json:"path,omitempty"`
	Descriptions       []string `json:"descriptions"`
	Aliases            []string `json:"aliases"`
	PrevExecutionNames []string `json:"prev"`
	Lang               string   `json:"lang"`
	Script             string   `json:"script"`
	Environments       []string `json:"environment"`
	WorkingDir         string   `json:"working_dir"`
	AdditionalArgs     []string `json:"additional_args"`
}

// ToJSONContent converts the instance to ExecutionJSONContent.
func (e *Execution) ToJSONContent() ExecutionJSONContent {

	prevNames := make([]string, 0, len(e.PrevExecutions))
	for _, prev := range e.PrevExecutions {
		prevNames = append(prevNames, prev.String())
	}

	return ExecutionJSONContent{
		Name:               e.Name.String(),
		Descriptions:       e.Descriptions,
		Path:               e.Path.String(),
		PrevExecutionNames: prevNames,
		Lang:               e.Lang.String(),
		Script:             e.Script,
		Environments:       e.Environments,
		WorkingDir:         e.WorkingDir,
		AdditionalArgs:     e.AdditionalArgs,
	}
}

// MarshalJSON represents the instance in JSON format.
func (e *Execution) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToJSONContent())
}
