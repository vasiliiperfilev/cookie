package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestPostConversation(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	conversationModel := data.NewStubConversationModel([]data.Conversation{})
	userModel := data.NewStubUserModel(generateUsers(4))
	models := data.Models{User: userModel, Conversation: conversationModel}
	t.Run("it POST conversation", func(t *testing.T) {
		server := app.New(cfg, logger, models)
		dto := data.PostConversationDto{
			UserIds: []int64{1, 2},
		}
		expectedResponse := data.Conversation{
			Id:            1,
			UserIds:       []int64{1, 2},
			LastMessageId: 0,
			Version:       1,
		}
		// post request
		request := createPostConversationRequest(t, dto)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		// assertion
		var gotConversation data.Conversation
		json.NewDecoder(response.Body).Decode(&gotConversation)
		assertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		data.AssertConversation(t, gotConversation, expectedResponse)
	})

	t.Run("it POST and GET same conversation by any of ids", func(t *testing.T) {
		server := app.New(cfg, logger, models)
		userIds := []int64{3, 4}
		userInput := data.PostConversationDto{
			UserIds: userIds,
		}
		expectedResponse := []data.Conversation{
			{
				Id:            2,
				UserIds:       userIds,
				LastMessageId: 0,
				Version:       1,
			},
		}
		// request
		postRequest := createPostConversationRequest(t, userInput)
		postRequest.Header.Set("Authorization", "Bearer "+strings.Repeat("3", 26))
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		// assertions
		for _, id := range userIds {
			getRequest := createGetAllConversationRequest(t, id)
			getRequest.Header.Set("Authorization", "Bearer "+strings.Repeat("3", 26))
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			var gotConversations []data.Conversation
			json.NewDecoder(response.Body).Decode(&gotConversations)

			assertStatus(t, response.Code, http.StatusOK)
			assertContentType(t, response, app.JsonContentType)
			data.AssertConversation(t, gotConversations[0], expectedResponse[0])
		}
	})

	t.Run("it don't allow POST same conversation", func(t *testing.T) {
		server := app.New(cfg, logger, models)
		userIds := []int64{3, 4}
		userInput := data.PostConversationDto{
			UserIds: userIds,
		}
		var response *httptest.ResponseRecorder
		for i := 0; i < 2; i++ {
			response = httptest.NewRecorder()
			postRequest := createPostConversationRequest(t, userInput)
			postRequest.Header.Set("Authorization", "Bearer "+strings.Repeat("3", 26))
			server.ServeHTTP(response, postRequest)
		}
		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
	})

	t.Run("can't PUT", func(t *testing.T) {
		server := app.New(cfg, logger, models)
		request, err := http.NewRequest(http.MethodPut, "/v1/conversations", nil)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
		assertHeader(t, response.Header().Get("Allow"), http.MethodPost, http.MethodGet)
	})
}

func createGetAllConversationRequest(t *testing.T, userId int64) *http.Request {
	getRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/conversations?userId=%v", userId), nil)
	tester.AssertNoError(t, err)
	return getRequest
}

func createPostConversationRequest(t *testing.T, dto data.PostConversationDto) *http.Request {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)
	request, err := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
	tester.AssertNoError(t, err)
	return request
}
