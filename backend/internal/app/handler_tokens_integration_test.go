package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestIntegrationTokenPost(t *testing.T) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		database.POSTGRES_USER,
		database.POSTGRES_PASSWORD,
		database.POSTGRES_PORT,
		database.POSTGRES_DB,
	)
	cfg := database.Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	db, err := database.OpenDB(cfg)
	tester.AssertNoError(t, err)

	t.Run("it returns a token after creating a user", func(t *testing.T) {
		server := app.PrepareIntegrationTestServer(db, 4000)
		email := "test5@nowhere.com"
		password := "test123!A"
		registerInput := data.RegisterUserInput{
			Email:    email,
			Password: password,
			Type:     1,
			ImageId:  "imageid",
		}
		userResponse := registerUser(t, server, registerInput)
		var user data.User
		json.NewDecoder(userResponse.Body).Decode(&user)

		loginInput := map[string]string{
			"Email":    email,
			"Password": password,
		}
		response := loginUser(t, server, loginInput)

		assertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		assertTokenResponse(t, response.Body, user.Id)
	})
	db.Close()
}

func loginUser(t *testing.T, server http.Handler, input map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)

	request, err := http.NewRequest(http.MethodPost, "/v1/tokens", requestBody)
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}