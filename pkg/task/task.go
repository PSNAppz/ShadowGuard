package task

import (
	"AegisGuard/pkg/monitor"
	"net/http"
)

type Type string

const (
	MONITOR Type = "monitor"
	CAPTURE Type = "capture"
)

type Task struct {
	Type     Type
	Settings map[string]interface{}
}

func (t *Task) Handle(r *http.Request) {
	switch t.Type {
	case MONITOR:
		monitor.Handle(r, t.Settings)
	default:
		panic("invalid task type given.")
	}
}
