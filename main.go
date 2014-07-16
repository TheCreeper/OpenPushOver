/*
    TODO:
        - Add a webs interface and message cache
*/

package main

import (

    "log"
    "time"
    "flag"
    "fmt"

    "github.com/TheCreeper/OpenPushOver/pushover"
    "github.com/TheCreeper/OpenPushOver/notification"
)

func (cfg *ClientConfig) launchListener(acn Account) {

    // Generate a UUID and save it to the config
    client := &pushover.Client{

        UserName: acn.Username,
        UserPassword: acn.Password,

        Key: acn.Key,

        DeviceName: cfg.Globals.DeviceName,
        DeviceUUID: acn.DeviceUUID,
    }

    // Specify some other options
    if (len(acn.Proxy) > 1) {

        conn := &ConnHandler {

            ProxyType: acn.proxyType,
            ProxyAddress: acn.proxyAddress,
            ProxyUsername: acn.proxyUsername,
            ProxyPassword: acn.proxyPassword,
            ProxyTimeout: acn.proxyTimeout,
        }
        client.Dial = conn.HandleConnection
    }

    err := client.LoginDevice()
    if (err != nil) {

        log.Fatalf("LoginDevice: %s\n", err)
    }

    if (acn.Register) {

        err = client.RegisterDevice(acn.Register)
        if (err != nil) {

            log.Fatal("RegisterDevice: %s\n", err)
        }
        acn.Register = false
        cfg.Flush()
    }

    for {

        time.Sleep(time.Duration(cfg.Globals.CheckFrequencySeconds) * time.Second)

        fetched, err := client.FetchMessages()
        if (err != nil) {

            log.Printf("FetchMessages: %s\n", err)
            continue

        }
        if (fetched > 0) {

            log.Printf("Fetched %d Messages\n", fetched)
        }

        for _, v := range client.MessagesResponse.Messages {

            var urgency string

            if (v.Priority == pushover.HighPriority) {

                urgency = notification.NormalPriority
            } else if (v.Priority == pushover.EmergencyPriority) {

                urgency = notification.CriticalPriority
            } else {

                urgency = notification.NormalPriority
            }

            if (len(v.Title) < 1) {

                v.Title = "Pushover Notification"
            }
            v.Title = fmt.Sprintf("%s (%s)", v.Title, time.Unix(v.Date, 0).Format("2006-01-02 15:04:05"))

            log.Print(v.Priority)

            n := &notification.Notify{

                Title: v.Title,
                Message: v.Message,
                Urgency: urgency,
                Icon: "dialog-information",
                Category: "im.received",
            }
            err = n.Push() // trigger the desktop notifications
            if (err != nil) {

                log.Printf("Push: %s\n", err)
            }
            log.Printf("[%d]: %s: %s\n", v.Id, v.Title, v.Message)
        }

        err = client.MarkRead()
        if (err != nil) {

            log.Fatalf("MarkRead: %s\n", err)
        }
    }
}

func init() {

    flag.StringVar(&configFile, "config", "./config.json", "The configuration file location")
    flag.Parse()
}

func main() {

    cfg, err := GetCFG(configFile)
    if (err != nil) {

        log.Fatalf("GetCFG: %s\n", err)
    }

    for i, v := range cfg.Accounts {

        if (len(v.DeviceUUID) < 1) {

            log.Print("Generating device UUID...")

            uuid, err := pushover.GenerateUUID()
            if (err != nil) {

                log.Fatalf("GenerateUUID: %s\n", err)
            }
            v.DeviceUUID = uuid
            cfg.Accounts[i].DeviceUUID = uuid

            err = cfg.Flush() // write changes to the config to disk
            if (err != nil) {

                log.Printf("cfg.Flush: %s\n", err)
            }

            log.Print("UUID Generated")
        }
        go cfg.launchListener(v)
    }

    select{}
}