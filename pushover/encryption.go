package pushover

import (

    "crypto/hmac"
    "crypto/sha256"
    "crypto/rand"
    "encoding/base64"
    "errors"

    "code.google.com/p/go.crypto/nacl/secretbox"
)

const (

    KeySize     = 32
    NonceSize   = 24
)

var (

    ErrSecretBox    = errors.New("Pushover: Failed to open Secretbox")
    ErrHMAC         = errors.New("Pushover: Unable to generate HMAC")
    ErrVerifyHMAC   = errors.New("Pushover: Unable to verify HMAC")
    ErrEncodeBase64 = errors.New("Pushover: Unable to encode to base64")

)

func (c *Client) DecryptMessage(message string) (out string, err error) {

    var key [KeySize]byte
    copy(key[:], c.Key)

    // Decode message
    decoded, err := DecodeBase64String(message)
    if (err != nil) {

        return
    }

    // Decrypt message
    decrypted, err := Decrypt(key, decoded)
    if (err != nil) {

        return
    }

    out = string(decrypted)
    return
}

func Decrypt(key [KeySize]byte, in []byte) (out []byte, err error) {

    var nonce [NonceSize]byte
    copy(nonce[:], in[:NonceSize])

    var ok bool
    out, ok = secretbox.Open(out, in[NonceSize:], &nonce, &key)
    if (!ok) {

        err = ErrSecretBox
        return
    }

    return
}

func (c *Client) EncryptMessage(message string) (out string, err error) {

    var key [KeySize]byte
    copy(key[:], c.Key)

    b, err := Encrypt(key, []byte(message))
    if (err != nil) {

        return
    }

    out, err = EncodeBase64String(b)
    if (err != nil) {

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

func EncodeBase64String(in []byte) (out string, err error) {

    out = base64.StdEncoding.EncodeToString(in)
    if (len(out) < 1) {

        err = ErrEncodeBase64
        return
    }

    return
}

func DecodeBase64String(in string) (out []byte, err error) {

    out, err = base64.StdEncoding.DecodeString(in)
    if (err != nil) {

        return
    }

    return
}

func GenerateHMAC(key, message []byte) (out []byte, err error) {

    h := hmac.New(sha256.New, key)
    h.Write(message)
    out = h.Sum(nil)
    if (out == nil) {

        err = ErrHMAC
        return
    }

    return
}

func VerifyHMAC(key, h, message []byte) (err error) {

    hmacExpected, err := GenerateHMAC(key, message)
    if (err != nil) {

        return
    }

    ok := hmac.Equal(h, hmacExpected)
    if !(ok) {

        err = ErrVerifyHMAC
        return
    }

    return
}