package main

import (
	"flag"
	"log"
	"time"

	"github.com/TheCreeper/OpenPushOver/pushover"
)

var (
	apptoken string
	userkey  string

	key string

	title    string
	message  string
	priority int
	expire   int
	sound    string
	url      string
	urltitle string
	callback string
)

func init() {

	flag.StringVar(&apptoken, "apptoken", "", "")
	flag.StringVar(&userkey, "userkey", "", "")
	flag.StringVar(&key, "key", "", "")

	flag.StringVar(&title, "title", "", "")
	flag.StringVar(&message, "message", "", "")
	flag.IntVar(&priority, "priority", 0, "")
	flag.IntVar(&expire, "expire", 15, "")
	flag.StringVar(&sound, "sound", "", "")
	flag.StringVar(&url, "url", "", "")
	flag.StringVar(&urltitle, "url-title", "", "")
	flag.StringVar(&callback, "callback", "", "")
	flag.Parse()
}

func main() {

	client := pushover.Client{

		AppToken: apptoken,
		UserKey:  userkey,
		Key:      key,
	}

	message := pushover.PushMessage{

		Title:     title,
		Message:   message,
		Priority:  priority,
		Expire:    expire,
		Url:       url,
		UrlTitle:  urltitle,
		Sound:     sound,
		Callback:  callback,
		Timestamp: int64(time.Now().Unix()),
	}

	var encrypt = false
	if len(key) > 1 {

		encrypt = true
	}

	err := client.PushMessage(message, encrypt)
	if err != nil {

		log.Fatalf("Push: %s\n", err)
	}

	log.Printf("Message Sent\n")
	log.Printf("AppLimit Messages: %s\n", client.Accounting.AppLimit)
	log.Printf("AppLimit Remaining Messages: %s\n", client.Accounting.AppRemaining)
	log.Printf("AppLimit Time to Reset: %s\n", client.Accounting.AppReset)
}
