package main

import (

    "flag"
    "log"

    "github.com/TheCreeper/Push/blob/master/pushover/"
)

var (

    apptoken string
    userkey string

    encryptionkey string

    title string
    message string
    priority int
    expire int
)

func init() {

    flag.StringVar(&apptoken, "apptoken", "", "")
    flag.StringVar(&userkey, "userkey", "", "")
    flag.StringVar(&title, "title", "", "")
    flag.StringVar(&message, "message", "", "")
    flag.IntVar(&priority, "priority", 0, "")
    flag.IntVar(&expire, "expire", 15, "")
    flag.StringVar(&encryptionkey, "key", "", "")
    flag.Parse()
}

func main() {

    client := pushover.Client{

        AppToken: apptoken,
        UserKey: userkey,
        Key: encryptionkey,
    }

    message := pushover.Message{

        Title: title,
        Message: message,
        Priority: priority,
        Expire: expire,
    }

    err := client.PushMessage(message)
    if (err != nil) {

        log.Fatalf("Push: %s\n", err)
    }
}