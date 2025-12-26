package providers

import (
	"fmt"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	err := &Error{Code: "PROVIDER_ERROR", Message: "boom"}
	if err.Error() != "boom" {
		t.Fatalf("expected boom, got %s", err.Error())
	}

	err.Err = fmt.Errorf("root")
	if err.Error() != "boom: root" {
		t.Fatalf("expected chained error, got %s", err.Error())
	}
}

func TestErrorMessageSelfReference(t *testing.T) {
	err := &Error{Code: "PROVIDER_ERROR", Message: "boom"}
	err.Err = err
	if err.Error() != "boom" {
		t.Fatalf("expected boom, got %s", err.Error())
	}
}
