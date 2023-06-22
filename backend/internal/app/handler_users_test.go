package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
	"golang.org/x/exp/slices"
)

// bad request if incorrect json
func TestUserPost(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{User: data.NewStubUserModel([]data.User{})}
	server := app.New(cfg, logger, models)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		userInput := data.PostUserDto{
			Email:    "test@nowhere.com",
			Name:     "test",
			Password: "test123!A",
			Type:     1,
			ImageId:  "imageid",
		}
		expectedResponse := data.User{
			Email:   userInput.Email,
			Name:    userInput.Name,
			Type:    userInput.Type,
			ImageId: userInput.ImageId,
		}
		requestBody := new(bytes.Buffer)
		request := createRegisterRequest(t, requestBody, userInput)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertRegisterResponse(t, response.Body, expectedResponse)
	})

	t.Run("it fails registration with duplicate email", func(t *testing.T) {
		userInput := data.PostUserDto{
			Email:    "test@nowhere.com",
			Password: "test123!A",
			Type:     1,
			ImageId:  "imageid",
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)

		request := createRegisterRequest(t, requestBody, userInput)
		response := httptest.NewRecorder()
		server.ServeHTTP(httptest.NewRecorder(), request)
		// second request with the same email
		request = createRegisterRequest(t, requestBody, userInput)
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
	})

	t.Run("can't POST with empty body", func(t *testing.T) {
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode("")
		request, _ := http.NewRequest(http.MethodPost, "/v1/users", requestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("can't PUT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/v1/users", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertHeader(t, response.Header().Get("Allow"), http.MethodPost, http.MethodGet)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestUserSearch(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	t.Run("it returns array of users with names similar to search query", func(t *testing.T) {
		r := "test+name"
		// generate 3 users with "test" in names
		users := generateUsers(4)
		// replace one user name with "name"
		users[2].Name = "name 1"
		// replace last user name with unmatching string
		users[3].Name = "unmatching"
		// setup a server
		models := data.Models{User: data.NewStubUserModel(users)}
		server := app.New(cfg, logger, models)
		// create and send request
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/users?q=%v", r), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		// expect
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		var got []data.User
		err = json.NewDecoder(response.Body).Decode(&got)
		tester.AssertNoError(t, err)
		// all users except unmatching
		want := users[:len(users)-1]
		assertUserSearch(t, got, want, r)
	})
}

func assertUserSearch(t *testing.T, got []data.User, want []data.User, r string) {
	if len(got) != len(want) {
		t.Fatalf("Expected to have %v responses", len(want))
	}
	keywords := strings.Split(r, "+")
	for _, usr := range got {
		nameKw := strings.Split(usr.Name, " ")
		fit := false
		for _, w := range nameKw {
			if slices.Contains(keywords, w) {
				fit = true
			}
		}
		if !fit {
			t.Fatalf("Expected name %v to match regexp %v", usr.Name, r)
		}
	}
}

func assertRegisterResponse(t *testing.T, body *bytes.Buffer, want data.User) {
	t.Helper()
	var got data.User
	json.NewDecoder(body).Decode(&got)
	assertUser(t, got, want)
}

func assertUser(t *testing.T, got data.User, want data.User) {
	t.Helper()
	if got.Email != want.Email {
		t.Fatalf("Expected email to be %v, got %v", want.Email, got.Email)
	}

	if got.Type != want.Type {
		t.Fatalf("Expected type to be %v, got %v", want.Type, got.Type)
	}

	if got.ImageId != want.ImageId {
		t.Fatalf("Expected email to be %v, got %v", want.ImageId, got.ImageId)
	}
}

func createRegisterRequest(t *testing.T, requestBody *bytes.Buffer, userInput data.PostUserDto) *http.Request {
	json.NewEncoder(requestBody).Encode(userInput)
	request, err := http.NewRequest(http.MethodPost, "/v1/users", requestBody)
	tester.AssertNoError(t, err)
	return request
}
