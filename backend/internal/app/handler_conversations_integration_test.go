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

func TestIntegrationConversations(t *testing.T) {
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

	t.Run("it posts and gets a conversation", func(t *testing.T) {
		server := app.PrepareIntegrationTestServer(db, 4000)
		email := "test55@nowhere.com"
		password := "test123!A"
		registerInput := data.PostUserDto{
			Email:    email,
			Password: password,
			Name:     "test",
			Type:     1,
			ImageId:  "imageid",
		}
		// register first user
		user1 := mustRegisterUser(t, server, registerInput)
		// register second user
		registerInput.Email = "test15@nowhere.com"
		user2 := mustRegisterUser(t, server, registerInput)
		// login as first user
		loginInput := map[string]string{
			"Email":    email,
			"Password": password,
		}
		userToken := mustLoginUser(t, server, loginInput)
		// post conversation
		dto := data.PostConversationDto{
			UserIds: []int64{user1.Id, user2.Id},
		}
		want := data.Conversation{
			Users:       []data.User{user1, user2},
			LastMessage: data.Message{Id: 0},
			Version:     1,
		}
		got := postConversation(t, server, userToken.Token.Plaintext, dto)
		want.Id = got.Id
		data.AssertConversation(t, got, want)
		// get conversations
		conversations := getConversations(t, server, userToken)
		for _, got := range conversations {
			if got.Id == want.Id {
				data.AssertConversation(t, got, want)
			}
		}
	})
}

func postConversation(t *testing.T, server http.Handler, token string, dto data.PostConversationDto) data.Conversation {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)

	request, err := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusCreated)
	var cvs data.Conversation
	json.NewDecoder(response.Body).Decode(&cvs)
	return cvs
}

func getConversations(t *testing.T, server *app.Application, token app.UserToken) []data.Conversation {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/conversations?userId=%v", token.User.Id), nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token.Plaintext))
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	var conversations []data.Conversation
	json.NewDecoder(response.Body).Decode(&conversations)

	return conversations
}
