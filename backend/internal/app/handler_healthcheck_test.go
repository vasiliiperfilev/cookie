package app_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestHealthcheckHandler(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{User: data.NewStubUserModel([]data.User{})}
	server := app.New(cfg, logger, models)
	t.Run("it returns health status", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)

		got := getAppStateFromResponse(t, response.Body)
		assertEnv(t, got.Env, "development")
	})
	t.Run("can't POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/v1/healthcheck", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertHeader(t, response.Header().Get("Allow"), http.MethodGet)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	tester.AssertValue(t, got, want, "Wrong http response status")
}

func assertHeader(t *testing.T, got string, want ...string) {
	gotArray := strings.Split(got, "; ")
	for _, wantVal := range want {
		contains := false
		for _, gotVal := range gotArray {
			if gotVal == wantVal {
				contains = true
				break
			}
		}
		if !contains {
			t.Fatalf("Header doesn't contain %v", wantVal)
		}
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	tester.AssertValue(t, response.Result().Header.Get("content-type"), want, "Wrong http response content-type")
}

func getAppStateFromResponse(t testing.TB, body io.Reader) (appState data.AppState) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&appState)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into env, '%v'", body, err)
	}

	return
}

func assertEnv(t *testing.T, got string, want string) {
	tester.AssertValue(t, got, want, "Wrong application Env config")
}
