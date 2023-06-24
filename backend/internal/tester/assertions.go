package tester

import (
	"reflect"
	"testing"
	"time"
)

func AssertValue[T any](t *testing.T, got T, want T, message string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%v, got %v, want %v", message, got, want)
	}
}

func AssertNoError(t *testing.T, err error) {
	AssertValue(t, err, nil, "Expected no error, but got one")
}

func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected to have a error, but got nil")
	}
}

func RetryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func AssertStatus(t *testing.T, got int, want int) {
	AssertValue(t, got, want, "Wrong http response status")
}
