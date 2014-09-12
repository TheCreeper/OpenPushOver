package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	MinCheckSeconds = 5
)

var (
	ConfigFile string
)

// Some errors
var (
	ErrNoDevName    = errors.New("No device name specified")
	ErrCheckSeconds = fmt.Errorf("No time specified for checkseconds or less than %s", MinCheckSeconds)
)

type ClientConfig struct {
	Globals Globals

	Proxys []struct {
		Name     string
		Type     string
		Address  string
		Username string
		Password string
		Timeout  int
	}

	Accounts []Account
}

type Globals struct {
	CacheDir     string
	DeviceName   string
	CheckSeconds int
}

type Account struct {
	DeviceUUID string

	Username string
	Password string

	Key string

	Proxy         string
	proxyType     string
	proxyAddress  string
	proxyUsername string
	proxyPassword string
	proxyTimeout  int
}

func (cfg *ClientConfig) Flush(f string) (err error) {

	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {

		return
	}
	defer file.Close()

	b, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {

		return
	}

	buf := bufio.NewWriter(file)

	_, err = buf.Write(b)
	if err != nil {

		return
	}
	defer buf.Flush()

	return
}

func (cfg *ClientConfig) validate() (err error) {

	if len(cfg.Globals.DeviceName) < 1 {

		return ErrNoDevName
	}

	if cfg.Globals.CheckSeconds < MinCheckSeconds {

		return ErrCheckSeconds
	}

	for i, v := range cfg.Accounts {

		if len(v.Proxy) < 1 {

			continue
		}

		for _, pv := range cfg.Proxys {

			if cfg.Accounts[i].Proxy == pv.Name {

				cfg.Accounts[i].proxyType = pv.Type
				cfg.Accounts[i].proxyAddress = pv.Address
				cfg.Accounts[i].proxyUsername = pv.Username
				cfg.Accounts[i].proxyPassword = pv.Password
				cfg.Accounts[i].proxyTimeout = pv.Timeout
			}
		}
	}

	return
}

func GetCFG(f string) (cfg *ClientConfig, err error) {

	b, err := ioutil.ReadFile(f)
	if err != nil {

		return
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {

		return
	}

	err = cfg.validate()
	if err != nil {

		return
	}

	return
}
