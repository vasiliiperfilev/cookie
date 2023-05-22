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
	server := app.New(cfg, logger)

	t.Run("it registers", func(t *testing.T) {
		userInput := data.RegisterUserInput{
			Email:    "test@nowhere.com",
			Password: "test123!A",
			Type:     "Supplier",
			ImageId:  "imageid",
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)
		request, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", requestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)

		userResponse := data.RegisterUserResponse{
			Email:   userInput.Email,
			Type:    userInput.Type,
			ImageId: userInput.ImageId,
		}
		// extract user from response
		var responseBody data.RegisterUserResponse
		json.NewDecoder(response.Body).Decode(&responseBody)

		if !reflect.DeepEqual(userResponse, responseBody) {
			t.Fatalf("Expected user to be %v, got %v", userResponse, responseBody)
		}
	})
}
