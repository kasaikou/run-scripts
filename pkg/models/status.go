package models

type Status int

const (
	StatusUndefined Status = iota
	StatusWaiting
	StatusRunning
	StatusSuccess
	StatusFailed
	StatusNotPlan
)

func (es Status) IsValid() bool {
	switch es {
	case StatusNotPlan, StatusWaiting, StatusRunning, StatusSuccess, StatusFailed:
		return true
	default:
		return false
	}
}

func (es Status) String() string {
	switch es {
	case StatusNotPlan:
		return "not_plan"
	case StatusWaiting:
		return "waiting"
	case StatusRunning:
		return "running"
	case StatusSuccess:
		return "success"
	case StatusFailed:
		return "failed"
	default:
		return "undefined"
	}
}
