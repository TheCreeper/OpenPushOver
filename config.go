package main

import (

    "encoding/json"
    "io/ioutil"
    "os"
    "bufio"
)

var (

    configFile string
)

type ClientConfig struct {

    Globals Globals

    Proxys []struct {

        Name string
        Type string
        Address string
        Username string
        Password string
        Timeout int
    }

    Accounts []Account

    Map struct {

        Sounds struct {

            Pushover string
            Bike string
            Bugle string
            Cashregister string
            Classical string
            Cosmic string
            Falling string
            Gamelan string
            Incomming string
            Intermission string
            Magic string
            Mechanical string
            Pianobar string
            Siren string
            Spacealarm string
            Tugboat string
            Alien string
            Climb string
            Persistent string
            Echo string
            Updown string
            None string
        }
    }
}

type Globals struct {

    DeviceName string
    CheckFrequencySeconds int
}

type Account struct {

    DeviceUUID string

    Register bool
    Username string
    Password string

    Key string

    Proxy string
    proxyType string
    proxyAddress string
    proxyUsername string
    proxyPassword string
    proxyTimeout int
}

func (cfg *ClientConfig) Flush() (err error) {

    file, err := os.OpenFile(configFile, os.O_RDWR | os.O_CREATE, 0666)
    if (err != nil) {

        return
    }
    defer file.Close()

    b, err := json.MarshalIndent(cfg, "", "    ")
    if (err != nil) {

        return
    }

    buf := bufio.NewWriter(file); defer buf.Flush()
    _, err = buf.Write(b)
    if (err != nil) {

        return
    }

    return
}

func (cfg *ClientConfig) validate() (err error) {

    if (len(cfg.Globals.DeviceName) < 1) {

        cfg.Globals.DeviceName = GetHostName()
    }

    if (cfg.Globals.CheckFrequencySeconds < 10) {

        cfg.Globals.CheckFrequencySeconds = 10 // Be friendly to their servers
    }

    for i, v := range cfg.Accounts {

        if (len(v.Proxy) < 1) {

            continue
        }

        for _, pv := range cfg.Proxys {

            if (cfg.Accounts[i].Proxy == pv.Name) {

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

func GetCFG(f string) (cfg ClientConfig, err error) {

    b, err := ioutil.ReadFile(f)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(b, &cfg)
    if (err != nil) {

        return
    }

    err = cfg.validate()
    if (err != nil) {

        return
    }

    return
}