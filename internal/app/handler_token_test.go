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
}
