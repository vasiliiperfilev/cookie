package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func TestPostToken(t *testing.T) {
	email := "test@test.com"
	password := "pa5$wOrd123"
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	user := data.User{
		Id:        1,
		CreatedAt: time.Now(),
		Email:     email,
		Type:      1,
		ImageId:   "id",
		Version:   1,
	}
	user.Password.Set(password)
	models := data.Models{User: data.NewStubUserModel([]data.User{user}), Token: data.NewStubTokenModel([]data.Token{})}
	server := app.New(cfg, logger, models)

	t.Run("it sends token response", func(t *testing.T) {
		userInput := struct {
			Email    string
			Password string
		}{
			Email:    email,
			Password: password,
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)
		request, _ := http.NewRequest(http.MethodPost, "/v1/token", requestBody)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		assertTokenResponse(t, response.Body, user.Id)
	})
}

func assertTokenResponse(t *testing.T, body *bytes.Buffer, userId int64) {
	t.Helper()
	var response app.HandlerTokenResponse
	v := validator.New()
	json.NewDecoder(body).Decode(&response)
	data.ValidateTokenPlaintext(v, response.Token.Plaintext)
	if len(v.Errors) != 0 {
		t.Fatalf("Incorrect token %v", v.Errors)
	}
	if response.User.Id != userId {
		t.Fatalf("Incorrect user returned, expected user id %v, got %v", userId, response.User.Id)
	}
}
