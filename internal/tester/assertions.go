package tester

import (
	"reflect"
	"testing"
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
		t.Fatalf("Expected to have a error %s, but got nil", err)
	}
}
