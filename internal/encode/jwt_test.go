package encode

import (
	"testing"
)

func TestGenerateJWT(t *testing.T) {

	header := Header{Alg: "SHA256", Typ: "JWT"}
	payload := Payload{Username: "test_user", Exp: 123123123}

	expectedToken := "eyJBbGciOiJTSEEyNTYiLCJUeXAiOiJKV1QifQ==." +
		"eyJVc2VybmFtZSI6InRlc3RfdXNlciIsIkV4cCI6MTIzMTIzMTIzfQ==." +
		"f9afc1e9ff93986c82548b68cd483712b6f5b8536bbc4e3e252dab3e2eb39eb5"

	token := GenerateJWT(header, payload, "test2")
	if token != expectedToken {
		t.Errorf("Expected '%s', got '%s'", expectedToken, token)
	}

}
