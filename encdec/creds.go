package encdec

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func decodeB64(str string) string {

	encodedCreds := strings.Split(str, " ")[1]
	rawDecodedText, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", rawDecodedText)
}

func GetCreds(str string) (string, string) {
	decoded := decodeB64(str)
	splitted := strings.Split(decoded, ":")
	cliendId := splitted[0]
	secret := splitted[1]
	return cliendId, secret
}
