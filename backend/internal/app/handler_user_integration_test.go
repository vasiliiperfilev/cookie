package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestIntegrationUserPost(t *testing.T) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		database.POSTGRES_USER,
		database.POSTGRES_PASSWORD,
		database.POSTGRES_PORT,
		database.POSTGRES_DB,
	)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		db := database.PrepareTestDb(t, dsn)
		server := app.PrepareServer(db, 4000)
		database.ApplyFixtures(t, db, "../fixtures")
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
		response := registerUser(t, server, userInput)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertRegisterResponse(t, response.Body, expectedResponse)
	})
}

func registerUser(t *testing.T, server http.Handler, input data.RegisterUserInput) *httptest.ResponseRecorder {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)

	request, err := http.NewRequest(http.MethodPost, "/v1/user", requestBody)
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	return response
}
