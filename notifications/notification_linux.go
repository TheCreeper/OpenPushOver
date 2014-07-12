package notifications

import (

    "os/exec"
    "errors"
    "fmt"
    "time"
)

func (n *Notify) Push() (err error) {

    if !(n.Timestamp < 1) {

        n.date = time.Unix(n.Timestamp, 0).Format("0:00")
    }

    if (len(n.Title) < 1) {

        n.Title = "PushOver Notification"
    }
    n.Title = fmt.Sprintf("%s (%d)", n.Title, n.date)

    if (len(n.Icon) < 1) {

        n.Icon = "dialog-information"
    }

    if (len(n.Urgency) < 1) {

        n.Urgency = "normal"
    }

    err = n.send()
    if (err != nil) {

        return
    }

    return
}

func (n *Notify) send() (err error) {

    cmd := exec.Command("notify-send", n.Title, n.Message, "--icon", n.Icon, "-u", n.Urgency)
    out, err := cmd.CombinedOutput();
    if (err != nil) {

        err = errors.New(string(out))
        return
    }

    return
}