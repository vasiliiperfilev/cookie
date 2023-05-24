package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

// can register
// can't not allowed methods
// can't register with the same email
// can't register without enough fields
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
		expectedResponse := data.RegisterUserResponse{
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
		assertResponseJson(t, response.Body, expectedResponse)
	})
}

func assertResponseJson[T any](t *testing.T, body *bytes.Buffer, expectedStruct T) {
	t.Helper()
	var response T
	json.NewDecoder(body).Decode(&response)

	if !reflect.DeepEqual(expectedStruct, response) {
		t.Fatalf("Expected user to be %v, got %v", expectedStruct, response)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
}
