package models

import (
	"fmt"

	"github.com/google/uuid"
)

// ExecutionReferer is interface for getting Execution associated with a name.
type ExecutionReferer interface {
	GetExecution(name string) *Execution
}

// ReferenceExecution is ExecutionReferer and associated name.
type ReferenceExecution struct {
	referer ExecutionReferer
	name    string
}

// NewReferenceExecution creates a ReferenceExecution instance.
func NewReferenceExecution(referer ExecutionReferer, name string) ReferenceExecution {
	return ReferenceExecution{referer: referer, name: name}
}

// Validate checks name is associated Execution in referer.
func (en ReferenceExecution) Validate() error {
	if execution := en.referer.GetExecution(en.name); execution == nil {
		return NewModelValidateError(fmt.Errorf("cannot found execution '%s'", en.name))
	}
	return nil
}

// String indicates the name of Execution associated with the name.
func (en ReferenceExecution) String() string {
	return en.referer.GetExecution(en.name).Name.String()
}

// ID indicates the ID of Execution associated with the name.
func (en ReferenceExecution) ID() uuid.UUID {
	return en.referer.GetExecution(en.name).ID
}
