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
		registerInput := data.PostUserDto{
			Email:    email,
			Password: password,
			Name:     "test",
			Type:     1,
			ImageId:  "imageid",
		}
		user := mustRegisterUser(t, server, registerInput)

		loginInput := map[string]string{
			"Email":    email,
			"Password": password,
		}
		got := mustLoginUser(t, server, loginInput)
		assertUserToken(t, got, user.Id)
	})
	db.Close()
}

func mustLoginUser(t *testing.T, server http.Handler, input map[string]string) app.UserToken {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)
	request, err := http.NewRequest(http.MethodPost, "/v1/tokens", requestBody)
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusCreated)
	var userToken app.UserToken
	json.NewDecoder(response.Body).Decode(&userToken)
	return userToken
}
