package encdec

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256SumHmacHex(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func Sha256SumB64(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
