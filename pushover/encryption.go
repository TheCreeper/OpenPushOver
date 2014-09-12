package pushover

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"code.google.com/p/go.crypto/nacl/secretbox"
)

// Encryption Limits
const (
	keySize   = 32
	nonceSize = 24
)

// Errors
var (
	ErrMsgNoEnc     = errors.New("Message is not encrypted")
	ErrSecretBox    = errors.New("Failed to open Secretbox")
	ErrHMAC         = errors.New("Unable to generate HMAC")
	ErrVerifyHMAC   = errors.New("Unable to verify HMAC")
	ErrEncodeBase64 = errors.New("Unable to encode to base64")
)

// Some regular expressions
var (

	// Valid encrypted message.
	ValidEncMsg = regexp.MustCompile("@Enc@.?")
)

func isEncrypted(msg string) bool {

	return ValidEncMsg.MatchString(msg)
}

func decryptMessage(key, s string) (msg string, err error) {

	var keyBytes [keySize]byte
	copy(keyBytes[:], key)

	s = ValidEncMsg.ReplaceAllString(s, "")

	// Decode message
	decoded, err := decodeBase64String(s)
	if err != nil {

		return
	}

	// Decrypt message
	decrypted, err := decrypt(keyBytes, decoded)
	if err != nil {

		return
	}

	msg = string(decrypted)
	return
}

func decrypt(key [keySize]byte, in []byte) (out []byte, err error) {

	var nonce [nonceSize]byte
	copy(nonce[:], in[:nonceSize])

	var ok bool
	out, ok = secretbox.Open(out, in[nonceSize:], &nonce, &key)
	if !ok {

		err = ErrSecretBox
		return
	}

	return
}

func (c *Client) encryptMessage(msg PushMessage) (err error) {

	var key [keySize]byte
	copy(key[:], c.Key)

	// Encrypt the message body
	b, err := encrypt(key, []byte(msg.Message))
	if err != nil {

		return
	}

	out, err := encodeBase64String(b)
	if err != nil {

		return
	}

	if len(out) > MessageLimit {

		return ErrMessageLimit
	}

	msg.Message = fmt.Sprintf("%s %s", "@Encrypted@", out)
	return
}

func encrypt(key [keySize]byte, in []byte) (out []byte, err error) {

	// Create a new nonce
	_, nonce, err := newNonce()
	if err != nil {

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
	if err != nil {

		return
	}

	return
}

func encodeBase64String(in []byte) (out string, err error) {

	out = base64.StdEncoding.EncodeToString(in)
	if len(out) < 1 {

		err = ErrEncodeBase64
		return
	}

	return
}

func decodeBase64String(in string) (out []byte, err error) {

	out, err = base64.StdEncoding.DecodeString(in)
	if err != nil {

		return
	}

	return
}
