package main

import (

    "os/exec"
    "errors"
    "fmt"
    "time"
)

func NotifySend(title, desc, icon, urgency string, timestamp int64) (err error) {

    var date string
    if (timestamp > 1) {

        date = time.Unix(timestamp, 0).Format("0:00")
    }

    if (len(title) < 1) {

        title = "PushOver Notification"
    }
    title = fmt.Sprintf("%s (%s)", title, date)

    if (len(icon) < 1) {

        icon = "dialog-information"
    }

    if (len(urgency) < 1) {

        urgency = "normal"
    }

    cmd := exec.Command("notify-send", title, desc, "--icon", icon, "-u", urgency)
    out, err := cmd.CombinedOutput();
    if (err != nil) {

        err = errors.New(string(out))
        return
    }

    return
}