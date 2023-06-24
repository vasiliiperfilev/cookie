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
		want.Id = got.Id
		want.Version = got.Version
		want.CreatedAt = got.CreatedAt
		data.AssertUser(t, got, want)
	})
	db.Close()
}

func TestIntegrationUserSearch(t *testing.T) {
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

	t.Run("it returns []users for search request", func(t *testing.T) {
		server := app.PrepareIntegrationTestServer(db, 4000)
		query := "veryspecialname"
		want := []data.User{}
		// register users with names
		for i := 1; i <= 5; i++ {
			userInput := data.PostUserDto{
				Email:    fmt.Sprintf("testReg%v@nowhere.com", i),
				Password: "test123!A",
				Name:     fmt.Sprintf("%v %v", query, i),
				Type:     1,
				ImageId:  "imageid",
			}
			user := mustRegisterUser(t, server, userInput)
			want = append(want, user)
		}
		// login as first user
		loginInput := map[string]string{
			"Email":    "testReg1@nowhere.com",
			"Password": "test123!A",
		}
		userToken := mustLoginUser(t, server, loginInput)
		got := mustGetAllUsersBySearch(query, t, userToken, server)
		assertUserSearch(t, got, want, query)
	})
	db.Close()
}

func mustGetAllUsersBySearch(query string, t *testing.T, userToken app.UserToken, server *app.Application) []data.User {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/users?q=%v", query), nil)
	tester.AssertNoError(t, err)
	request.Header.Set("Authorization", "Bearer "+userToken.Token.Plaintext)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	assertContentType(t, response, app.JsonContentType)
	var got []data.User
	err = json.NewDecoder(response.Body).Decode(&got)
	tester.AssertNoError(t, err)
	return got
}

func mustRegisterUser(t *testing.T, server http.Handler, input data.PostUserDto) data.User {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)

	request, err := http.NewRequest(http.MethodPost, "/v1/users", requestBody)
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	var user data.User
	json.NewDecoder(response.Body).Decode(&user)

	return user
}
