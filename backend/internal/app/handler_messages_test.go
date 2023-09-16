package app_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestMessagesHandler(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	userModel := data.NewStubUserModel(generateUsers(9))
	conversationModel := data.NewStubConversationModel(generateConversation(4), userModel)
	want := data.Message{Id: 1, ConversationId: 1, Content: "test", SenderId: 1, PrevMessageId: 0}
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{
		want,
	})
	models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
	app := app.New(cfg, logger, models)
	t.Run("it GET messages by conversation id", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/conversations/1/messages", nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusOK)
		var got []data.Message
		json.NewDecoder(response.Body).Decode(&got)
		// 0 message is sentinel node
		if !reflect.DeepEqual(got[1], want) {
			t.Fatalf("Want message %v, but got %v", want, got[0])
		}
	})

	t.Run("it 404 if user is not in conversation", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/conversations/5/messages", nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("9", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
