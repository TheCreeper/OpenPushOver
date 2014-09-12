package notification

import (
	"os/exec"
	"path/filepath"
)

func (m *Message) Push() (err error) {

	if (len(m.Title) < 1) && (len(m.Body) < 1) {

		return ErrTitleMsg
	}

	var args []string

	if len(m.Title) > 1 {

		args = append(args, m.Title)
	}

	if len(m.Body) > 1 {

		args = append(args, m.Body)
	}

	// A custom image needs to be an absolute path
	if len(m.Icon) > 1 {

		args = append(args, "--icon="+filepath.Clean(m.Icon))
	}

	if len(m.Urgency) > 1 {

		args = append(args, "--urgency="+m.Urgency)
	}

	if m.ExpireTime > 1 {

		args = append(args, "--expire-time="+string(m.ExpireTime))
	}

	if len(m.Category) > 1 {

		args = append(args, "--category="+m.Category)
	}

	if len(m.Hint) > 1 {

		args = append(args, "--hint="+m.Hint)
	}

	cmd := exec.Command("notify-send", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {

		return &NotificationErr{Return: string(out), Err: err}
	}

	if len(m.Sound) > 1 {

		return m.PlaySound()
	}

	return
}

func (m *Message) PlaySound() (err error) {

	cmd := exec.Command("paplay", m.Sound)
	out, err := cmd.CombinedOutput()
	if err != nil {

		return &NotificationErr{File: m.Sound, Return: string(out), Err: err}
	}

	return
}
