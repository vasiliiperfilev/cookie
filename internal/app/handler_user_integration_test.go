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
)

func TestIntegrationUserPost(t *testing.T) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DB)
	db := prepareTestDb(t, dsn)
	server := prepareServer(db, 4000)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		applyFixtures(t, db, "../fixtures")
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
	assertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	return response
}
