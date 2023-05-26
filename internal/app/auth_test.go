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
)

// bad request if incorrect json
func TestAuthRegister(t *testing.T) {
	env := "testing"
	cfg := app.Config{Port: 4000, Env: env}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{User: data.NewStubUserModel()}
	server := app.New(cfg, logger, models)

	t.Run("it returns correct response", func(t *testing.T) {
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
		json.NewEncoder(requestBody).Encode(userInput)

		request, err := http.NewRequest(http.MethodPost, "/v1/auth/register", requestBody)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertNoError(t, err)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertRegisterResponse(t, response.Body, expectedResponse)
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

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
}
