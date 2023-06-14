package encode

import (
	"testing"
)

func TestGetCreds(t *testing.T) {
	username, password := GetCreds("Basic dGVzdDp0ZXN0cGFzcw==")
	if username != "test" {
		t.Errorf("Expected 'expected result', got '%s'", username)
	}
	if password != "testpass" {
		t.Errorf("Expected 'expected result', got '%s'", password)
	}
}
