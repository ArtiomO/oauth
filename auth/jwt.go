package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Header struct {
	Alg string
	Typ string
}

type Payload struct {
	Username string
	Exp      int
}

type Stringer interface {
	String() string
}

func (c Payload) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (c Header) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func encodeToB64(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func sha256Sum(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func structToString(stringer Stringer) string {
	return stringer.String()
}

func GenerateJWT(h Header, p Payload, secret string) string {
	headerEncoded := encodeToB64(structToString(h))
	payloadEncoded := encodeToB64(structToString(p))
	payloadWithHeader := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	shaSum := sha256Sum(payloadWithHeader, secret)
	return fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, shaSum)
}
