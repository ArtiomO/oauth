package encode

import (
	"testing"
)

func TestSha256SumHmacHex(t *testing.T) {
	encrypted := Sha256SumHmacHex("testesttest", "testsecret")
	if encrypted != "ab0df45a886a6d6f2dd9b363f9d0f437e64a6fb57e4321e959b2d5dc23d8ae3f" {
		t.Errorf("Expected 'expected result', got '%s'", encrypted)
	}
}

func TestSha256SumB64(t *testing.T) {
	encrypted := Sha256SumHex("testtesttest")
	if encrypted != "a2c96d518f1099a3b6afe29e443340f9f5fdf1289853fc034908444f2bcb8982" {
		t.Errorf("Expected 'expected result', got '%s'", encrypted)
	}
}
