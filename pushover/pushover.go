/*
    Reference Implementation: github.com/Mechazawa/Pushover-client-protocol
    Reference Implementation: github.com/nbrownus/pushover-desktop-client
    Reference Implementation: github.com/AlekSi/pushover
    API Specs: pushover.net/api

    TODO:
        - Improve encryption with public/private keys
        - Account for errors returned by the API
        - Fix message priority not being parsed by fetchmessages
        - Add disable device
        - Do better validiation
        - Fetch client images
        - Return more detailed errors
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
    "strconv"
)

const (

    BaseUrl = "https://api.pushover.net/1"
)


// Message priority.
const (

    LowestPriority      = -2 // lowest priority, no notification
    LowPriority         = -1 // low priority, no sound and vibration
    NormalPriority      = 0 // normal priority, default
    HighPriority        = 1 // high priority, always with sound and vibration
    EmergencyPriority   = 2 // emergency priority, requires acknowledge
)

// Message sound.
const (

    PushoverSound       = "pushover"
    BikeSound           = "bike"
    BugleSound          = "bugle"
    CashregisterSound   = "cashregister"
    ClassicalSound      = "classical"
    CosmicSound         = "cosmic"
    FallingSound        = "falling"
    GamelanSound        = "gamelan"
    IncomingSound       = "incoming"
    IntermissionSound   = "intermission"
    MagicSound          = "magic"
    MechanicalSound     = "mechanical"
    PianobarSound       = "pianobar"
    SirenSound          = "siren"
    SpacealarmSound     = "spacealarm"
    TugboatSound        = "tugboat"
    AlienSound          = "alien"
    ClimbSound          = "climb"
    PersistentSound     = "persistent"
    EchoSound           = "echo"
    UpdownSound         = "updown"
    NoneSound           = "none"
)

var (

    ErrLoginFailed  = errors.New("Pushover: Failed to login")
    ErrDeviceReg    = errors.New("Pushover: Device register failed")
    ErrFetchMsgF    = errors.New("Pushover: Message fetch failed")
    ErrMarkReadF    = errors.New("Pushover: Markread messages failed")
    ErrPushMsgF     = errors.New("Pushover: Unable to push message")
    ErrReceiptF     = errors.New("Pushover: Unable to get receipt")
    ErrDeviceAuth   = errors.New("Pushover: Device not authenticated")
    ErrUserPassword = errors.New("Pushover: User Password not specified")
    ErrUserName     = errors.New("Pushover: UserName not specified")
    ErrMessageLimit = errors.New("Pushover: Message is longer than the limit")
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
    Errors []string `json:"errors"`
}

func (c *Client) LoginDevice() (err error) {

    // Some validiation
    if (len(c.UserName) < 1) {

        return ErrUserName
    }
    if (len(c.UserPassword) < 1) {

        return ErrUserPassword
    }

    err = VerifyDeviceName(c.DeviceName)
    if (err != nil) {

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

        return ErrLoginFailed
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
    Errors []string `json:"errors"`
}

func (c *Client) RegisterDevice(replaceDevice bool) (err error) {

    if (len(c.Login.Secret) < 1) {

        return ErrDeviceAuth
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

        return ErrDeviceReg
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
    Errors []string `json:"errors"`
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

        err = ErrDeviceAuth
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

        err = ErrFetchMsgF
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
    fetched = len(c.MessagesResponse.Messages)

    // Decrypt any encrypted messages
    if (len(c.Key) > 1) {

        err = c.DecryptMessages()
        if (err != nil) {

            return
        }
    }

    return
}

func (c *Client) DecryptMessages() (err error) {

    re, err := regexp.CompilePOSIX("@Encrypted@.?")
    if (err != nil) {

        return
    }

    for i, v := range c.MessagesResponse.Messages {

        if !(re.MatchString(v.Message)) {

            continue
        }
        v.Message = re.ReplaceAllString(v.Message, "")

        out, err := c.DecryptMessage(v.Message)
        if (err != nil) {

            c.MessagesResponse.Messages[i].Message = err.Error()
            continue
        }
        c.MessagesResponse.Messages[i].Message = out
    }

    return
}

type MarkReadResponse struct {

    Status int `json: "status"`
    Request string `json: "request"`
    Errors []string `json:"errors"`
}

func (c *Client) MarkRead() (err error) {

    if (len(c.Login.Secret) < 1) {

        return ErrDeviceAuth
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

        return ErrMarkReadF
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

    Receipt string `json:"receipt"`

    Status int `json: "status"`
    Request string `json: "request"`
    Errors []string `json:"errors"`
}

func (c *Client) Push(message string) (err error) {

    msg := Message {

        Message: message,
    }
    return c.PushMessage(msg, false)
}

type Message struct {

    Device      string // Device name to send to specific devices instead of all
    Title       string // Title of message
    Message     string // Message to send
    Priority    int    // Defaults to 0 although expire needs to be specified when higher than 1
    Url         string // Url to be sent with message
    UrlTitle    string // Url title
    Timestamp   int64 // Timestamp which should be a unixstamp
    Sound       string // Sound to be played on client device

    // Emergency notifications
    Expire      int
    Retry       int
    Callback    string
}

func (c *Client) PushMessage(msg Message, encrypt bool) (err error) {

    // Some validations
    err = VerifyUserKey(c.AppToken)
    if (err != nil) {

        return
    }

    err = VerifyUserKey(c.UserKey)
    if (err != nil) {

        return
    }

    err = VerifyPushMessage(msg)
    if (err != nil) {

        return
    }

    if (encrypt) {

        out, err := c.EncryptMessage(msg.Message)
        if (err != nil) {

            return err
        }
        if (len(out) > MessageLimit) {

            return ErrMessageLimit
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
    vars.Add("timestamp", strconv.FormatInt(msg.Timestamp, 10))
    vars.Add("sound", msg.Sound)
    vars.Add("callback", msg.Callback)

    urlF := fmt.Sprintf("%s%s", BaseUrl, "/messages.json")
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.PostForm(urlF, vars)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        return ErrPushMsgF
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
    c.Accounting.AppReset = resp.Header.Get("X-Limit-App-Reset")

    return
}

type Receipt struct {

    Acknowledged int `json:"acknowledged"`
    AcknowledgedAt int `json:"acknowledged_at"`
    LastDeliveredAt int `json:"last_delivered_at"`
    Expired int `json:"expired"`
    ExpiresAt int64 `json:"expires_at"`
    CalledBack int `json:"called_back"`
    CalledBackAt int64 `json:"called_back_at"`

    Status int `json:"status"`
    Request string `json: "request"`
    Errors []string `json:"errors"`
}

func (c *Client) GetReceipt(receipt string) (err error) {

    err = VerifyReceipt(receipt)
    if (err != nil) {

        return
    }

    urlF := fmt.Sprintf("%s/%s.json?token=%s", BaseUrl, "/receipts", receipt, c.AppToken)
    httpClient := &http.Client{Transport: &http.Transport{Dial: c.dial}}
    resp, err := httpClient.Get(urlF)
    if (err != nil) {

        return
    }
    if (resp.StatusCode >= 400) {

        return ErrReceiptF
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

    return
}
