package app_test

import (
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

func TestAuthenticateRequest(t *testing.T) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		database.POSTGRES_USER,
		database.POSTGRES_PASSWORD,
		database.POSTGRES_PORT,
		database.POSTGRES_DB,
	)
	db := database.PrepareTestDb(t, dsn)
	server := app.PrepareServer(db, 4000)

	t.Run("it returns a user from token", func(t *testing.T) {
		database.ApplyFixtures(t, db, "../fixtures")
		email := "test@nowhere.com"
		password := "test123!A"
		registerInput := data.RegisterUserInput{
			Email:    email,
			Password: password,
			Type:     1,
			ImageId:  "imageid",
		}
		registerUser(t, server, registerInput)

		loginInput := map[string]string{
			"Email":    email,
			"Password": password,
		}
		response := loginUser(t, server, loginInput)

		var token data.Token
		json.NewDecoder(response.Body).Decode(&token)

		user := authRequest(t, server, token)
		tester.AssertValue(t, user.Email, registerInput.Email, "same email")
		tester.AssertValue(t, user.ImageId, registerInput.ImageId, "same image id")
		tester.AssertValue(t, user.Type, registerInput.Type, "same type")
	})
}

func authRequest(t *testing.T, server *app.Application, token data.Token) *data.User {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, "/v1/token", nil)
	tester.AssertNoError(t, err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Plaintext))
	response := httptest.NewRecorder()
	user, err := server.AuthenticateRequest(response, request)
	tester.AssertNoError(t, err)
	return user
}
