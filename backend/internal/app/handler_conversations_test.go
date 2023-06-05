package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestPostConversation(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{Conversation: data.NewStubConversationModel([]*data.Conversation{})}
	server := app.New(cfg, logger, models)
	t.Run("it POST conversation", func(t *testing.T) {
		models := data.Models{Conversation: data.NewStubConversationModel([]*data.Conversation{})}
		server := app.New(cfg, logger, models)
		userInput := data.Conversation{
			UserIds: []int64{1, 2},
		}
		expectedResponse := data.Conversation{
			Id:            1,
			UserIds:       []int64{1, 2},
			LastMessageId: -1,
		}
		// post request
		request := createPostConversationRequest(t, userInput)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		// assertion
		var gotConversation data.Conversation
		json.NewDecoder(response.Body).Decode(&gotConversation)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertConversation(t, gotConversation, expectedResponse)
	})

	t.Run("it POST and GET same conversation by any of ids", func(t *testing.T) {
		models := data.Models{Conversation: data.NewStubConversationModel([]*data.Conversation{})}
		server := app.New(cfg, logger, models)
		userIds := []int64{3, 4}
		userInput := data.Conversation{
			UserIds: userIds,
		}
		expectedResponse := []data.Conversation{
			{
				Id:            1,
				UserIds:       userIds,
				LastMessageId: -1,
			},
		}
		// request
		postRequest := createPostConversationRequest(t, userInput)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		// assertions
		for _, id := range userIds {
			getRequest := createGetAllConversationRequest(t, id)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			var gotConversations []data.Conversation
			json.NewDecoder(response.Body).Decode(&gotConversations)

			assertStatus(t, response.Code, http.StatusOK)
			assertContentType(t, response, app.JsonContentType)
			assertConversation(t, gotConversations[0], expectedResponse[0])
		}
	})

	t.Run("it don't allow POST same conversation", func(t *testing.T) {
		models := data.Models{Conversation: data.NewStubConversationModel([]*data.Conversation{})}
		server := app.New(cfg, logger, models)
		userIds := []int64{3, 4}
		userInput := data.Conversation{
			UserIds: userIds,
		}
		var response *httptest.ResponseRecorder
		for i := 0; i < 2; i++ {
			response = httptest.NewRecorder()
			postRequest := createPostConversationRequest(t, userInput)
			server.ServeHTTP(response, postRequest)
		}
		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
	})

	t.Run("can't PUT", func(t *testing.T) {
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

func createPostConversationRequest(t *testing.T, userInput data.Conversation) *http.Request {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(userInput)
	request, err := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
	tester.AssertNoError(t, err)
	return request
}

func assertConversation(t *testing.T, got data.Conversation, want data.Conversation) {
	t.Helper()
	tester.AssertValue(t, want.Id, got.Id, "Expected same conversation id")
	tester.AssertValue(t, want.Id, got.Id, "Expected same last message id")
	sort.Slice(got.UserIds, func(i, j int) bool {
		return got.UserIds[i] >= got.UserIds[j]
	})
	sort.Slice(want.UserIds, func(i, j int) bool {
		return want.UserIds[i] >= want.UserIds[j]
	})
	tester.AssertValue(t, want.UserIds, got.UserIds, "Expected same usersId")
}
