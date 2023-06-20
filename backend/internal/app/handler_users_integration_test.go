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
	cfg := database.Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	db, err := database.OpenDB(cfg)
	tester.AssertNoError(t, err)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		server := app.PrepareIntegrationTestServer(db, 4000)
		userInput := data.PostUserDto{
			Email:    "testReg@nowhere.com",
			Password: "test123!A",
			Name:     "test",
			Type:     1,
			ImageId:  "imageid",
		}
		want := data.User{
			Email:   userInput.Email,
			Type:    userInput.Type,
			Name:    userInput.Name,
			ImageId: userInput.ImageId,
		}
		got := mustRegisterUser(t, server, userInput)
		data.AssertUser(t, got, want)
	})
	db.Close()
}

func mustRegisterUser(t *testing.T, server http.Handler, input data.PostUserDto) data.User {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)

	request, err := http.NewRequest(http.MethodPost, "/v1/users", requestBody)
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusOK)
	var user data.User
	json.NewDecoder(response.Body).Decode(&user)

	return user
}
