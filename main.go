/*
   TODO:
       - Add a webs interface and message cache
       - Allow custom sounds
*/

package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TheCreeper/OpenPushOver/notification"
	"github.com/TheCreeper/OpenPushOver/pushover"
)

// Create a new logger
var log = logrus.New()

// Map pushover prioritys to Notify prioritys
var PushoverToNotifyPriority = map[int]string{

	pushover.LowestPriority:  notification.LowPriority,
	pushover.LowPriority:     notification.LowPriority,
	pushover.NormalPriority:  notification.NormalPriority,
	pushover.HighPriority:    notification.CriticalPriority,
	pushover.HighestPriority: notification.CriticalPriority,
}

func (cfg *ClientConfig) launchPushover(acn Account) {

	// Generate a UUID and save it to the config
	client := &pushover.Client{

		UserName:     acn.Username,
		UserPassword: acn.Password,

		Key: acn.Key,

		DeviceName: cfg.Globals.DeviceName,
		DeviceUUID: acn.DeviceUUID,
	}

	// Specify some other options
	if len(acn.Proxy) > 1 {

		conn := &ConnHandler{

			ProxyType:     acn.proxyType,
			ProxyAddress:  acn.proxyAddress,
			ProxyUsername: acn.proxyUsername,
			ProxyPassword: acn.proxyPassword,
			ProxyTimeout:  acn.proxyTimeout,
		}
		client.Dial = conn.HandleConnection
	}

	err := client.LoginDevice()
	if err != nil {

		log.Errorf("LoginDevice: %s", err)
		return
	}

	if acn.Register {

		err = client.RegisterDevice(acn.Register)
		if err != nil {

			log.Errorf("RegisterDevice: %s", err)
			return
		}
		acn.Register = false
		cfg.Flush(ConfigFile)
	}

	for {

		time.Sleep(time.Duration(cfg.Globals.CheckFrequencySeconds) * time.Second)

		fetched, err := client.FetchMessages()
		if err != nil {

			log.Warn(err)
			continue

		}
		if fetched > 0 {

			log.Infof("Fetched %d Messages", fetched)
		}

		for _, v := range client.MessagesResponse.Messages {

			// Check if quiet hours is enabled
			if (client.MessagesResponse.User.QuietHours) && (v.Priority == pushover.NormalPriority) {

				v.Priority = pushover.LowPriority
			}

			var snd string
			// Check if sound file exists
			if len(v.Sound) > 1 {

				f, err := filepath.Abs(filepath.Join(cfg.Globals.CacheDir, pushover.SoundFileName[v.Sound]))
				if err != nil {

					log.Error(err)
					return
				}
				snd = f

				exists, err := FileExists(snd)
				if err != nil {

					log.Warn(err)
				}
				if !exists {

					b, err := client.FetchSound(v.Sound)
					if err != nil {

						log.Warn(err)
					}

					err = WriteToFile(snd, b)
					if err != nil {

						log.Warn(err)
					}
				}
			}

			// Make sure it has the default icon
			// TODO: Allow custom default image
			if len(v.Icon) < 1 {

				v.Icon = "default"
			}

			var img string
			// Check if image file exists
			if len(v.Icon) > 1 {

				f, err := filepath.Abs(filepath.Join(cfg.Globals.CacheDir, fmt.Sprintf("%s.png", v.Icon)))
				if err != nil {

					log.Error(err)
					return
				}
				img = f

				exists, err := FileExists(img)
				if err != nil {

					log.Warn(err)
				}
				if !exists {

					b, err := client.FetchImage(v.Icon)
					if err != nil {

						log.Warn(err)
					}

					err = WriteToFile(img, b)
					if err != nil {

						log.Warn(err)
					}
				}
			}

			// Make sure it has title
			if len(v.Title) < 1 {

				v.Title = "Pushover Notification"
			}
			v.Title = fmt.Sprintf("%s (%s)", v.Title, time.Unix(v.Date, 0).Format("2006-01-02 15:04:05"))

			// trigger the desktop notifications
			n := &notification.Notify{

				Title:    v.Title,
				Message:  v.Message,
				Urgency:  PushoverToNotifyPriority[v.Priority],
				Icon:     img,
				Category: "im.received",
				Sound:    snd,
			}
			err = n.Push()
			if err != nil {

				log.Warn(err)
			}

			// Print the notification to terminal
			log.Infof("[%d]: %s: %s", v.Id, v.Title, v.Message)
		}

		err = client.MarkRead()
		if err != nil {

			log.Warn(err)
		}
	}
}

func init() {

	flag.StringVar(&ConfigFile, "config", "./config.json", "The configuration file location")
	flag.Parse()
}

func main() {

	cfg, err := GetCFG(ConfigFile)
	if err != nil {

		log.Errorf("GetCFG: %s", err)
		return
	}

	for i, v := range cfg.Accounts {

		if len(v.DeviceUUID) < 1 {

			log.Info("Generating device UUID...")

			uuid, err := pushover.GenerateUUID()
			if err != nil {

				log.Errorf("GenerateUUID: %s", err)
				return
			}
			v.DeviceUUID = uuid
			cfg.Accounts[i].DeviceUUID = uuid

			err = cfg.Flush(ConfigFile) // write changes to the config to disk
			if err != nil {

				log.Warnf("cfg.Flush: %s", err)
			}
		}
		go cfg.launchPushover(v)
	}

	select {}
}
