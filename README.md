OpenPushOver Client
=====================

- Requirements

    - GNU/Linux only
    - libnotify (GNU/Linux) for notifications


##Sample Config

```json
{
    "Globals": {
        "DeviceName": "Fusion",
        "CheckFrequencySeconds": 10
    },
    "Proxys": [
        {
            "Name": "Tor",
            "Type": "socks5",
            "Address": "127.0.0.1:9050",
            "Username": "",
            "Password": "",
            "Timeout": 1
        }
    ],
    "Accounts": [
        {
            "DeviceUUID": "",
            "Register": true,
            "Username": "email",
            "Password": "password",
            "Key": "testkey123456789",
            "Proxy": "Tor"
        }
    ]
}
```