package encode

import (
	b64 "encoding/base64"
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

func structToString(stringer Stringer) string {
	return stringer.String()
}

func GenerateJWT(h Header, p Payload, secret string) string {
	headerEncoded := encodeToB64(structToString(h))
	payloadEncoded := encodeToB64(structToString(p))
	payloadWithHeader := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	shaSum := Sha256SumHmacHex(payloadWithHeader, secret)
	return fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, shaSum)
}
