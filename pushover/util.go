package pushover

import(

    "runtime"
    "regexp"
    "fmt"
    "errors"

    "github.com/nu7hatch/gouuid"
)

const (

    UserKeyLimit        = 30
    AppTokenLimit       = 30

    MessageTitleLimit   = 100
    MessageLimit        = 512

    UrlTitleLimit       = 100
    UrlLimit            = 512

    DeviceNameLimit     = 25
    DeviceUUIDLimit     = 36

    ReceiptLimit        = 30
)

var (

    ErrDeviceName   = errors.New("Pushover: DeviceName must contain at least one character and may only contain letters, numbers, dashes, and underscores")
    ErrUserKey      = fmt.Errorf("Pushover: User and group identifiers must be at least %d characters long, case-sensitive, and may only contain letters and numbers\n", UserKeyLimit)
    ErrAppToken     = fmt.Errorf("Pushover: Application tokens are case-sensitive and must be at least %d characters long, and may only contain letters and numbers\n", AppTokenLimit)
    ErrDeviceUUID   = fmt.Errorf("Pushover: Device UUID must be no longer than %d characters long", DeviceUUIDLimit)
    ErrReceipt      = fmt.Errorf("Pushover: Receipt must be at least %s and is case-sensitive", ReceiptLimit)

    ErrMsgLimit     = fmt.Errorf("Pushover: Message specified is not specified or is over the %d char limit\n", MessageLimit)
    ErrTitleLimit   = fmt.Errorf("Pushover: Title specified is over the %d char limit\n", MessageTitleLimit)
    ErrUrlTLimit    = fmt.Errorf("Pushover: Url Title specified is over the %d char limit\n", UrlTitleLimit)
    ErrUrlLimit     = fmt.Errorf("Pushover: The url is over the %d char limit\n", UrlLimit)
    ErrPriority     = errors.New("Pushover: A priority higher than 1 needs an expiry parm")
)

func GenerateUUID() (id string, err error) {

    u4, err := uuid.NewV4()
    if err != nil {

        return
    }

    id = u4.String()
    return
}

func GetHostOS() (os string) {

    if (runtime.GOOS != "") {

        return runtime.GOOS
    }

    return "unknown"
}

func VerifyDeviceName(name string) (err error) {

    if (len(name) < 1) {

        return
    }

    if (len(name) > DeviceNameLimit) {

        return ErrDeviceName
    }

    re, err := regexp.CompilePOSIX("^[a-zA-Z0-9_-]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(name)) {

        return ErrDeviceName
    }

    return
}

func VerifyUserKey(key string) (err error) {

    if (len(key) < 1) || (len(key) > UserKeyLimit) {

        return ErrUserKey
    }

    re, err := regexp.CompilePOSIX("^[A-Za-z0-9]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(key)) {

        return ErrUserKey
    }

    return
}

func VerifyAppToken(token string) (err error) {

    if (len(token) < 1) || (len(token) > AppTokenLimit) {

        return ErrAppToken
    }

    re, err := regexp.CompilePOSIX("^[A-Za-z0-9]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(token)) {

        return ErrAppToken
    }

    return
}

func VerifyDeviceUUID(id string) (err error) {

    if (len(id) < 1) {

        return ErrDeviceUUID
    }

    _, err = uuid.Parse([]byte(id))
    if (err != nil) {

        return
    }

    // Use regexp too to verify

    return
}

func VerifyPushMessage(msg Message) (err error) {

    if (len(msg.Message) < 1) || (len(msg.Message) > MessageLimit) {

        return ErrMsgLimit
    }
    if (len(msg.Title) > MessageTitleLimit) {

        return ErrTitleLimit
    }
    if (len(msg.UrlTitle) > UrlTitleLimit) {

        return ErrUrlTLimit
    }
    if (len(msg.Url) > UrlLimit) {

        return ErrUrlLimit
    }
    if (msg.Priority > 1 && msg.Expire < 1) {

        return ErrPriority
    }

    return
}

func VerifyReceipt(receipt string) (err error) {

    if (len(receipt) < 1) || (len(receipt) > ReceiptLimit) {

        return ErrReceipt
    }

    re, err := regexp.CompilePOSIX("^[A-Za-z0-9]{30}$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(receipt)) {

        return ErrReceipt
    }

    return
}