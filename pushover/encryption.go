package pushover

import (

    "crypto/rand"
    "encoding/base64"
    "errors"

    "code.google.com/p/go.crypto/nacl/secretbox"
)

const (

    KeySize = 32
    NonceSize = 24
)

func (c *Client) DecryptMessage(message string) (out string, err error) {

    var key [KeySize]byte
    copy(key[:], c.Key)

    // Decode message
    b, err := DecodeBase64String(message)
    if (err != nil) {

        return
    }

    // Decrypt message
    out, err = Decrypt(key, b)
    if (err != nil) {

        return
    }

    return
}

func Decrypt(key [KeySize]byte, in []byte) (out string, err error) {

    var nonce [NonceSize]byte
    copy(nonce[:], in[:NonceSize])

    var ok bool
    var b []byte
    b, ok = secretbox.Open(b, in[NonceSize:], &nonce, &key)
    if (!ok) {

        err = errors.New("Failed to open Secretbox")
        return
    }
    out = string(b)

    return
}

func (c *Client) EncryptMessage(message string) (out string, err error) {

    var key [KeySize]byte
    copy(key[:], c.Key)

    b, err := Encrypt(key, []byte(message))
    if (err != nil) {

        return
    }

    out = EncodeBase64String(b)
    if (len(out) < 1) {

        return
    }

    return
}

func Encrypt(key [KeySize]byte, in []byte) (out []byte, err error) {

    // Create a new nonce
    _, nonce, err := NewNonce()
    if (err != nil) {

        return
    }

    // Encrypt
    out = secretbox.Seal(out, []byte(in), &nonce, &key)

    // Put the nonce at the front of the array
    out = append(nonce[:], out...)

    return
}

func NewNonce() (i int, nonce [NonceSize]byte, err error) {

    i, err = rand.Read(nonce[:])
    if (err != nil) {

        return
    }

    return
}

func EncodeBase64String(s []byte) string {

    return base64.StdEncoding.EncodeToString(s)
}

func DecodeBase64String(s string) (d []byte, err error) {

    d, err = base64.StdEncoding.DecodeString(s)
    if (err != nil) {

        return
    }

    return
}