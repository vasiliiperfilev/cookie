package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

// bad request if incorrect json
func TestUserPost(t *testing.T) {
	env := "testing"
	cfg := app.Config{Port: 4000, Env: env}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{User: data.NewStubUserModel([]data.User{})}
	server := app.New(cfg, logger, models)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		userInput := data.RegisterUserInput{
			Email:    "test@nowhere.com",
			Password: "test123!A",
			Type:     1,
			ImageId:  "imageid",
		}
		expectedResponse := data.User{
			Email:   userInput.Email,
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
		userInput := data.RegisterUserInput{
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
		request, _ := http.NewRequest(http.MethodPost, "/v1/user", requestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("can't GET", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/user", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertHeader(t, response.Header().Get("Allow"), http.MethodPost)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func assertRegisterResponse(t *testing.T, body *bytes.Buffer, expectedUser data.User) {
	t.Helper()
	var response data.User
	json.NewDecoder(body).Decode(&response)

	if response.Email != expectedUser.Email {
		t.Fatalf("Expected email to be %v, got %v", expectedUser.Email, response.Email)
	}

	if response.Type != expectedUser.Type {
		t.Fatalf("Expected type to be %v, got %v", expectedUser.Type, response.Type)
	}

	if response.ImageId != expectedUser.ImageId {
		t.Fatalf("Expected email to be %v, got %v", expectedUser.ImageId, response.ImageId)
	}
}

func createRegisterRequest(t *testing.T, requestBody *bytes.Buffer, userInput data.RegisterUserInput) *http.Request {
	json.NewEncoder(requestBody).Encode(userInput)
	request, err := http.NewRequest(http.MethodPost, "/v1/user", requestBody)
	tester.AssertNoError(t, err)
	return request
}
