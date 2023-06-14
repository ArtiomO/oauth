package encode

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
)

type header struct {
	Alg string
	Typ string
}

type payload struct {
	Username string
	Exp      int
}

func (c payload) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (c header) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func encodeToB64(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func structToString(stringer fmt.Stringer) string {
	return stringer.String()
}

func headerFactory(alg, typ string) header {
	return header{Alg: alg, Typ: typ}
}

func payloadFactory(username string, exp int) payload {
	return payload{Username: username, Exp: exp}
}

func GenerateJWT(username string) string {

	header := headerFactory("SHA256", "JWT")
	payload := payloadFactory(username, 123123123)
	pay1 := structToString(header)
	pay2 := structToString(payload)
	headerEncoded := encodeToB64(pay1)
	payloadEncoded := encodeToB64(pay2)
	payloadWithHeader := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	shaSum := Sha256SumHmacHex(payloadWithHeader, "test2")
	return fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, shaSum)
}
