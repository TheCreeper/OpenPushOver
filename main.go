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
	"sync"
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

func (cfg *ClientConfig) LaunchClient(wg *sync.WaitGroup, acn *Account) {

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
	defer wg.Done()

	if len(acn.DeviceUUID) < 1 {

		err = client.RegisterDevice()
		if err != nil {

			log.Errorf("RegisterDevice: %s", err)
			return
		}
		acn.DeviceUUID = client.DeviceUUID

		err = cfg.Flush(ConfigFile)
		if err != nil {

			log.Errorf("Flush: %s", err)
			return
		}
	}

	for {

		time.Sleep(time.Duration(cfg.Globals.CheckSeconds) * time.Second)

		fetched, err := client.FetchMessages()
		if err != nil {

			log.Warn(err)
			continue

		}
		if fetched > 0 {

			log.Infof("Fetched %d Messages", fetched)

			for _, v := range client.MessagesResponse.Messages {

				// Check if quiet hours is enabled
				if (client.MessagesResponse.User.QuietHours) && (v.Priority == pushover.NormalPriority) {

					v.Priority = pushover.LowPriority
				}

				var snd string
				// Check if sound file exists
				if len(v.Sound) > 1 {

					f, err := filepath.Abs(filepath.Join(cfg.Globals.CacheDir, v.Sound+".wav"))
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

				v.Title = fmt.Sprintf("%s (%s)", v.Title, time.Unix(v.Date, 0).Format("2006-01-02 15:04:05"))

				// trigger the desktop notifications
				n := &notification.Message{

					Title:    v.Title,
					Body:     v.Message,
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
				log.Infof("[%d]: %s: %s", v.ID, v.Title, v.Message)
			}

			err = client.MarkReadHighest()
			if err != nil {

				log.Warn(err)
			}
		}
	}

	wg.Done()
}

func init() {

	flag.StringVar(&ConfigFile, "config", "./config.json", "The configuration file location")
	flag.Parse()
}

func main() {

	var wg sync.WaitGroup

	cfg, err := GetCFG(ConfigFile)
	if err != nil {

		log.Errorf("GetCFG: %s", err)
		return
	}

	for i := range cfg.Accounts {

		v := &cfg.Accounts[i]

		wg.Add(1)
		go cfg.LaunchClient(&wg, v)
	}

	wg.Wait()
}
