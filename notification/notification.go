package notification

import (
	"errors"
	"strings"
)

type Notify struct {
	Title      string
	Message    string
	Icon       string
	Urgency    string
	ExpireTime int
	Category   string
	Hint       string
	Sound      string
}

const (
	LowPriority      = "low"
	NormalPriority   = "normal"
	CriticalPriority = "critical"
)

// Errors
var (
	ErrTitleMsg = errors.New("Notification: A title or message must be specified")
)

type NotifyErr struct {
	File   string
	Return string
	Err    error
}

func (e *NotifyErr) Error() string {

	// Usually return will have a newline character
	return e.File + " " + strings.TrimSpace(e.Return) + ": " + e.Err.Error()
}
