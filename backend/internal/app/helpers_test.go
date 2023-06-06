package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestReadJson(t *testing.T) {
	t.Run("it errors if extra keys", func(t *testing.T) {
		got := struct {
			TestField  string
			ExtraField string
		}{
			TestField:  "test",
			ExtraField: "extra",
		}
		want := struct {
			TestField string
		}{
			TestField: "",
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(got)
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", requestBody)
		err := readJsonFromBody(httptest.NewRecorder(), request, &want)
		tester.AssertError(t, err)
	})

	t.Run("it errors if incorrectJson", func(t *testing.T) {
		requestBody := []byte(`{"TestField" - "test"}`)
		want := struct {
			TestField string
		}{
			TestField: "",
		}
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		err := readJsonFromBody(httptest.NewRecorder(), request, &want)
		tester.AssertError(t, err)
	})

	t.Run("it errors if incorrect field type", func(t *testing.T) {
		requestBody := []byte(`{"TestField":"test"}`)
		want := struct {
			TestField int
		}{
			TestField: 0,
		}
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		err := readJsonFromBody(httptest.NewRecorder(), request, &want)
		tester.AssertError(t, err)
	})

	t.Run("it errors if body is empty", func(t *testing.T) {
		requestBody := []byte(``)
		want := struct {
			TestField int
		}{
			TestField: 0,
		}
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		err := readJsonFromBody(httptest.NewRecorder(), request, &want)
		tester.AssertError(t, err)
	})

	t.Run("it errors if there is something after the JSON", func(t *testing.T) {
		requestBody := []byte(`{"TestField":"test"}asd`)
		want := struct {
			TestField string
		}{
			TestField: "",
		}
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		err := readJsonFromBody(httptest.NewRecorder(), request, &want)
		tester.AssertError(t, err)
	})
}
