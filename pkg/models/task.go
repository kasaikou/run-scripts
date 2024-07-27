package models

import (
	"context"

	"github.com/google/uuid"
)

type Task interface {
	Kind() string
	ID() uuid.UUID
	Start(ctx context.Context, req StartEventRequest)
	Depends() []Task
	IsPrepared() bool
}

type StartEventRequest struct {
	OnStart   func()
	OnSucceed func()
	OnFail    func()
}

type TaskNameRepository struct {
	TaskFromName(task Task, name []string)
}
