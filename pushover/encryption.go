package pushover

import (

    "crypto/rand"
    "encoding/base64"
    "errors"

    "code.google.com/p/go.crypto/nacl/secretbox"
)

const (

    keySize     = 32
    nonceSize   = 24
)

var (

    ErrSecretBox    = errors.New("Pushover: Failed to open Secretbox")
    ErrHMAC         = errors.New("Pushover: Unable to generate HMAC")
    ErrVerifyHMAC   = errors.New("Pushover: Unable to verify HMAC")
    ErrEncodeBase64 = errors.New("Pushover: Unable to encode to base64")

)

func (c *Client) DecryptMessage(message string) (out string, err error) {

    var key [keySize]byte
    copy(key[:], c.Key)

    // Decode message
    decoded, err := decodeBase64String(message)
    if (err != nil) {

        return
    }

    // Decrypt message
    decrypted, err := decrypt(key, decoded)
    if (err != nil) {

        return
    }

    out = string(decrypted)
    return
}

func decrypt(key [keySize]byte, in []byte) (out []byte, err error) {

    var nonce [nonceSize]byte
    copy(nonce[:], in[:nonceSize])

    var ok bool
    out, ok = secretbox.Open(out, in[nonceSize:], &nonce, &key)
    if (!ok) {

        err = ErrSecretBox
        return
    }

    return
}

func (c *Client) EncryptMessage(message string) (out string, err error) {

    var key [keySize]byte
    copy(key[:], c.Key)

    b, err := encrypt(key, []byte(message))
    if (err != nil) {

        return
    }

    out, err = encodeBase64String(b)
    if (err != nil) {

        return
    }

    return
}

func encrypt(key [keySize]byte, in []byte) (out []byte, err error) {

    // Create a new nonce
    _, nonce, err := newNonce()
    if (err != nil) {

        return
    }

    // Encrypt
    out = secretbox.Seal(out, []byte(in), &nonce, &key)

    // Put the nonce at the front of the array
    out = append(nonce[:], out...)

    return
}

func newNonce() (i int, nonce [nonceSize]byte, err error) {

    i, err = rand.Read(nonce[:])
    if (err != nil) {

        return
    }

    return
}

func encodeBase64String(in []byte) (out string, err error) {

    out = base64.StdEncoding.EncodeToString(in)
    if (len(out) < 1) {

        err = ErrEncodeBase64
        return
    }

    return
}

func decodeBase64String(in string) (out []byte, err error) {

    out, err = base64.StdEncoding.DecodeString(in)
    if (err != nil) {

        return
    }

    return
}