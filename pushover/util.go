package pushover

import(

    "runtime"
    "regexp"
    "errors"

    "github.com/nu7hatch/gouuid"
)

const (

    UserKeyLimit = 30
    AppTokenLimit = 30
    MessageTitleLimit = 100
    MessageLimit = 512
    UrlTitleLimit = 100
    UrlLimit = 512
    DeviceNameLimit = 25
)

var errorMessages = []string{

    "DeviceName must contain at least one character and may only contain letters, numbers, dashes, and underscores",
    "User and group identifiers must be at least 30 characters long, case-sensitive, and may only contain letters and numbers",
    "Application tokens are case-sensitive and must be at least 30 characters long, and may only contain letters and numbers",
}

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

    if (len(name) < 1) || (len(name) > DeviceNameLimit) {

        err = errors.New(errorMessages[0])
        return
    }

    re, err := regexp.CompilePOSIX("^[a-zA-Z0-9_-]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(name)) {

        err = errors.New(errorMessages[0])
        return
    }

    return
}

func VerifyUserKey(key string) (err error) {

    if (len(key) < 1) || (len(key) > UserKeyLimit) {

        err = errors.New(errorMessages[1])
        return
    }

    re, err := regexp.CompilePOSIX("^[A-Za-z0-9]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(key)) {

        err = errors.New(errorMessages[1])
        return
    }

    return
}

func VerifyAppToken(token string) (err error) {

    if (len(token) < 1) || (len(token) > AppTokenLimit) {

        err = errors.New(errorMessages[2])
        return
    }

    re, err := regexp.CompilePOSIX("^[A-Za-z0-9]+$")
    if (err != nil) {

        return
    }

    if !(re.MatchString(token)) {

        err = errors.New(errorMessages[2])
        return
    }

    return
}

func VerifyUUID(id string) (err error) {

    _, err = uuid.Parse([]byte(id))
    if (err != nil) {

        return
    }

    // Use regexp too to verify

    return
}