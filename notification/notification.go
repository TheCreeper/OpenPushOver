package notification

import (
	"errors"
	"strings"
)

type Message struct {
	Title      string
	Body       string
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

type NotificationErr struct {
	File   string
	Return string
	Err    error
}

func (e *NotificationErr) Error() string {

	// Usually return will have a newline character
	return e.File + " " + strings.TrimSpace(e.Return) + ": " + e.Err.Error()
}
