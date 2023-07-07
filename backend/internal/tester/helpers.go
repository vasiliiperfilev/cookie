package tester

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func ParseResponse[T any](t *testing.T, response *httptest.ResponseRecorder) T {
	t.Helper()
	var got T
	err := json.NewDecoder(response.Body).Decode(&got)
	AssertNoError(t, err)
	return got
}
