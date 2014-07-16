package notification

import (

    "os/exec"
    "errors"
)

var (

    errDesc = errors.New("Notifications: A title or message must be specified")
)

func (n *Notify) Push() (err error) {

    if (len(n.Title) < 1) && (len(n.Message) < 1) {

        return errDesc
    }

    var args []string

    if (len(n.Title) > 1) {

        args = append(args, n.Title)
    }

    if (len(n.Message) > 1) {

        args = append(args, n.Message)
    }

    if (len(n.Icon) > 1) {

        args = append(args, "--icon=" + n.Icon)
    }

    if (len(n.Urgency) > 1) {

        args = append(args, "--urgency=" + n.Urgency)
    }

    if (n.ExpireTime > 1) {

        args = append(args, "--expire-time=" + string(n.ExpireTime))
    }

    if (len(n.Category) > 1) {

        args = append(args, "--category=" + n.Category)
    }

    if (len(n.Hint) > 1) {

        args = append(args, "--hint=" + n.Hint)
    }

    cmd := exec.Command("notify-send", args...)
    out, err := cmd.CombinedOutput()
    if (err != nil) {

        err = errors.New(string(out))
        return
    }

    return
}