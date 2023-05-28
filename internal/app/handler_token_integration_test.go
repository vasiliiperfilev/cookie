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
)

// can register and log in
func TestIntegrationTokenPost(t *testing.T) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DB)
	db := prepareTestDb(t, dsn)
	server := prepareServer(db, 4000)

	t.Run("it returns a token after creating a user", func(t *testing.T) {
		applyFixtures(t, db, "../fixtures")
		email := "test@nowhere.com"
		password := "test123!A"
		registerInput := data.RegisterUserInput{
			Email:    email,
			Password: password,
			Type:     1,
			ImageId:  "imageid",
		}
		registerUser(t, server, registerInput)

		loginInput := struct {
			Email    string
			Password string
		}{
			email,
			password,
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(loginInput)

		request, err := http.NewRequest(http.MethodPost, "/v1/token", requestBody)
		assertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		assertTokenResponse(t, response.Body)
	})
}
