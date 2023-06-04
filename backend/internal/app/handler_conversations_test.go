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
			SupplierId: 1,
			BusinessId: 2,
		}
		// expectedResponse := data.Conversation{
		// 	ConversationId: 1,
		// 	SupplierId:     1,
		// 	BusinessId:     2,
		// 	LastMessageId:  -1,
		// }
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)
		request, err := http.NewRequest(http.MethodPost, "/v1/conversations", requestBody)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
	})
}
