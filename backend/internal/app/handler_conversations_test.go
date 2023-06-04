package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestPostConversation(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.Models{}
	server := app.New(cfg, logger, models)
	t.Run("it POST conversation", func(t *testing.T) {
		userInput := data.Conversation{
			UserIds: []int64{1, 2},
		}
		expectedResponse := data.Conversation{
			Id:            1,
			UserIds:       []int64{1, 2},
			LastMessageId: -1,
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)
		request, err := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		var gotConversation data.Conversation
		json.NewDecoder(response.Body).Decode(&gotConversation)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, gotConversation, expectedResponse, "expected different conversation")
	})

	t.Run("it GET conversation", func(t *testing.T) {
		expectedResponse := []data.Conversation{
			{
				Id:            1,
				UserIds:       []int64{1, 2},
				LastMessageId: -1,
			},
		}
		request, err := http.NewRequest(http.MethodGet, "/v1/conversations?userId=1", nil)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		var gotConversations []data.Conversation
		json.NewDecoder(response.Body).Decode(&gotConversations)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, gotConversations, expectedResponse, "expected different conversation")
	})

	t.Run("it POST and GET same conversation", func(t *testing.T) {
		userInput := data.Conversation{
			UserIds: []int64{4, 5},
		}
		expectedResponse := []data.Conversation{
			{
				Id:            1,
				UserIds:       []int64{4, 5},
				LastMessageId: -1,
			},
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)
		postRequest, _ := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)

		getRequest, _ := http.NewRequest(http.MethodGet, "/v1/conversations?userId=1", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getRequest)

		var gotConversations []data.Conversation
		json.NewDecoder(response.Body).Decode(&gotConversations)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, gotConversations, expectedResponse, "expected different conversation")
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
