/*
    Reference Implementation: github.com/Mechazawa/Pushover-client-protocol
    API Specs: pushover.net/api

    TODO:
        - Improve encryption with public/private keys
        - Account for errors returned by the API
        - Fix message priority not being parsed by fetchmessages
        - Add disable device
        - Do better validiation
*/

package pushover

import (

    "net"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "errors"
    "fmt"
    "regexp"
)

const (

    BaseUrl = "https://api.pushover.net/1"
)

type Client struct {

    Dial func(network, addr string) (net.Conn, error)

    UserName string
    UserPassword string

    DeviceName string
    DeviceUUID string
    deviceOS string
    provider_device_id string

    Key string

    Login Login
    Device Device
    MessagesResponse MessagesResponse
    MarkReadResponse MarkReadResponse

    AppToken string
    UserKey string

    Accounting struct {

        AppLimit string
        AppRemaining string
        AppReset string
    }

    PushResponse PushResponse
}

func (c *Client) dial(network, addr string) (net.Conn, error) {

    if (c.Dial != nil) {

        return c.Dial(network, addr)
    }

    dialer := &net.Dialer {

        DualStack: true,
    }

    return dialer.Dial(network, addr)
}

type Login struct {

    Status int `json:"status"`
    Secret string `json:"secret"`
    Request string `json:"request"`
    Id string `json:"id"`
}

func (c *Client) LoginDevice() (err error) {

    // Some validiation
    if (len(c.UserName) < 1) {

        err = errors.New("No username specified")
        return
    }
    if (len(c.UserPassword) < 1) {

        err = errors.New("No password specified")
        return
    }
    if (len(c.DeviceName) < 1) {

        err = errors.New("No devicename specified")
        return
    }
    if (len(c.DeviceName) > DeviceNameLimit) {

        err = errors.New("Device name is over the 25 char limit")
        return
    }
    if (len(c.DeviceUUID) < 1) {

        err = errors.New("No device uuid specified")
        return
    }

    // Set the unexported feilds
    c.deviceOS = GetHostOS()
    c.provider_device_id = c.deviceOS

    vars := url.Values{}
    vars.Add("email", c.UserName)
    vars.Add("password", c.UserPassword)

    urlF := fmt.Sprintf("%s%s", BaseUrl, "/users/login.json")
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.PostForm(urlF, vars)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        err = errors.New("Login Failed!")
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(body, &c.Login)
    if (err != nil) {

        return
    }

    return
}

type Device struct {

    Status int `json:"status"`
    Request string `json:"request"`
    Id string `json:"id"`
}

func (c *Client) RegisterDevice(replaceDevice bool) (err error) {

    if (len(c.Login.Secret) < 1) {

        err = errors.New("Device not authenticated!")
        return
    }

    err = VerifyDeviceName(c.DeviceName)
    if (err != nil) {

        return
    }

    var force = "0"
    if (replaceDevice) {

        force = "1"
    }

    vars := url.Values{}
    vars.Add("secret", c.Login.Secret)
    vars.Add("name", c.DeviceName)
    vars.Add("uuid", c.DeviceUUID)
    vars.Add("on_gcm", "1")
    vars.Add("os", c.deviceOS)
    vars.Add("force", force)
    vars.Add("provider_device_id", c.provider_device_id)

    urlF := fmt.Sprintf("%s%s", BaseUrl, "/devices.json")
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.PostForm(urlF, vars)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        err = errors.New("Device Register Failed")
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(body, &c.Device)
    if (err != nil) {

        return
    }

    return
}

type MessagesResponse struct {

    Messages []Messages `json:"messages"`

    User User `json:"user"`

    Status int `json:"status"`
    Request string `json:"request"`
}

type Messages struct {

    Id int `json: "id"`
    Message string `json: "message"`
    App string `json: "app"`
    Aid int `json: "aid"`
    Icon int `json: "icon"`
    Date int64 `json: "date"`
    Priority int `json: "priority"`
    Acked int `json: "acked"`
    Umid int `json: "umid"`
    Title string `json: "title"`
}

type User struct {

    QuietHours bool `json: "quiet_hours"`
    IsAndroid bool `json: "is_android_licensed"`
    IsIOS bool `json: "is_ios_licensed"`
    isDesktop bool `json: "is_desktop_licensed"`
}

func (c *Client) FetchMessages() (fetched int, err error) {

    if (len(c.Login.Secret) < 1) {

        err = errors.New("Device not authenticated!")
        return
    }

    vars := url.Values{}
    vars.Add("secret", c.Login.Secret)
    vars.Add("device_id", c.DeviceUUID)

    urlF := fmt.Sprintf("%s%s%s", BaseUrl, "/messages.json?", vars.Encode())
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.Get(urlF)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        err = errors.New("Message Fetch Failed")
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(body, &c.MessagesResponse)
    if (err != nil) {

        return
    }

    re, err := regexp.CompilePOSIX("@Encrypted@.?")
    if (err != nil) {

        return
    }

    for i, v := range c.MessagesResponse.Messages {

        if (len(c.Key) < 1) {

            continue
        }

        if !(re.MatchString(v.Message)) {

            continue
        }
        v.Message = re.ReplaceAllString(v.Message, "")

        out, err := c.DecryptMessage(v.Message)
        if (err != nil) {

            continue
        }
        c.MessagesResponse.Messages[i].Message = out
    }
    fetched = len(c.MessagesResponse.Messages)

    return
}

type MarkReadResponse struct {

    Status int `json: "status"`
    Request string `json: "request"`
}

func (c *Client) MarkRead() (err error) {

    if (len(c.Login.Secret) < 1) {

        err = errors.New("Device not authenticated!")
        return
    }

    vars := url.Values{}
    vars.Add("secret", c.Login.Secret)
    vars.Add("message", "99999")

    urlF := fmt.Sprintf("%s%s%s%s", BaseUrl, "/devices/", c.DeviceUUID,"/update_highest_message.json")
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.PostForm(urlF, vars)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        err = errors.New("MarkRead Failed")
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(body, &c.MarkReadResponse)
    if (err != nil) {

        return
    }

    return
}

type PushResponse struct {

    Status int
    Request string
}

func (c *Client) Push(message string) (err error) {

    msg := Message {

        Message: message,
    }
    return c.PushMessage(msg)
}

type Message struct {

    Device string
    Title string
    Message string
    Expire int
    Retry int
    Priority int
    Url string
    UrlTitle string
    Timestamp string
    Sound string
}

func (c *Client) PushMessage(msg Message) (err error) {

    // Some validations
    err = VerifyUserKey(c.AppToken)
    if (err != nil) {

        return
    }

    err = VerifyUserKey(c.UserKey)
    if (err != nil) {

        return
    }

    if (len(c.DeviceName) > 1) {

        err = VerifyDeviceName(c.DeviceName)
        if (err != nil) {

            return
        }
    }

    if (len(msg.Message) < 1) || (len(msg.Message) > MessageLimit) {

        err = errors.New("Message must be specified or is over the 512 char limit")
        return
    }
    if (len(msg.Title) > MessageTitleLimit) {

        err = errors.New("The title is over the 100 char limit")
        return
    }
    if (len(msg.Url) > UrlLimit) {

        err = errors.New("The url is over the 512 char limit")
        return
    }
    if (len(msg.UrlTitle) > UrlTitleLimit) {

        err = errors.New("The url title is over the 100 char limit")
        return
    }
    if (msg.Priority > 1 && msg.Expire < 1) {

        err = errors.New("A prioty of higher than 1 needs an expiry param")
        return
    }

    if (len(c.Key) > 1) {

        out, err := c.EncryptMessage(msg.Message)
        if (err != nil) {

            return err
        }
        if (len(out) > MessageLimit) {

            return errors.New("Encrypted Message is longer than the limit")
        }
        msg.Message = fmt.Sprintf("%s %s", "@Encrypted@", out)
    }

    vars := url.Values{}
    // Required
    vars.Add("token", c.AppToken)
    vars.Add("user", c.UserKey)
    vars.Add("message", msg.Message)
    // Optional
    vars.Add("device", msg.Device)
    vars.Add("title", msg.Title)
    vars.Add("url", msg.Url)
    vars.Add("url_title", msg.UrlTitle)
    vars.Add("expire", string(msg.Expire))
    vars.Add("retry", string(msg.Retry))
    vars.Add("priority", string(msg.Priority))
    vars.Add("timestamp", msg.Timestamp)
    vars.Add("sound", msg.Sound)

    urlF := fmt.Sprintf("%s%s", BaseUrl, "/messages.json")
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.PostForm(urlF, vars)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        err = errors.New("Push Message Failed")
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {

        return
    }

    err = json.Unmarshal(body, &c.PushResponse)
    if (err != nil) {

        return
    }

    // Update accounting
    c.Accounting.AppLimit = resp.Header.Get("X-Limit-App-Limit")
    c.Accounting.AppRemaining = resp.Header.Get("X-Limit-App-Remaining")
    c.Accounting.AppReset = resp.Header.Get("-Limit-App-Reset")

    return
}