package notification

import (
	"os/exec"
	"path/filepath"
)

func (n *Notify) Push() (err error) {

	if (len(n.Title) < 1) && (len(n.Message) < 1) {

		return ErrTitleMsg
	}

	var args []string

	if len(n.Title) > 1 {

		args = append(args, n.Title)
	}

	if len(n.Message) > 1 {

		args = append(args, n.Message)
	}

	// A custom image needs to be an absolute path
	if len(n.Icon) > 1 {

		args = append(args, "--icon=" + filepath.Clean(n.Icon))
	}

	if len(n.Urgency) > 1 {

		args = append(args, "--urgency=" + n.Urgency)
	}

	if n.ExpireTime > 1 {

		args = append(args, "--expire-time=" + string(n.ExpireTime))
	}

	if len(n.Category) > 1 {

		args = append(args, "--category=" + n.Category)
	}

	if len(n.Hint) > 1 {

		args = append(args, "--hint=" + n.Hint)
	}

	cmd := exec.Command("notify-send", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {

		return &NotifyErr{Return: string(out), Err: err}
	}

	if len(n.Sound) > 1 {

		return n.PlaySound()
	}

	return
}

func (n *Notify) PlaySound() error {

	cmd := exec.Command("paplay", n.Sound)
	out, err := cmd.CombinedOutput()
	if err != nil {

		return &NotifyErr{File: n.Sound, Return: string(out), Err: err}
	}

	return nil
}