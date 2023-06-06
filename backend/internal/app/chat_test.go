package app_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestChat(t *testing.T) {
	env := "testing"
	cfg := app.Config{Port: 4000, Env: env}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	messageModel := data.NewStubMessageModel([]data.Conversation{{Id: 1, UserIds: []int64{1, 2}}}, []data.Message{})
	models := data.Models{Message: messageModel}
	appServer := app.New(cfg, logger, models)
	t.Run("establishes ws connection", func(t *testing.T) {
		server := httptest.NewServer(appServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/v1/chat"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		want := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test",
		}

		js, err := json.Marshal(want)
		tester.AssertNoError(t, err)

		if err := ws.WriteMessage(websocket.TextMessage, js); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		messages, err := messageModel.GetAllByUserId(1)
		tester.AssertNoError(t, err)
		if len(messages) == 0 {
			t.Fatalf("Expected messages len to be > 0")
		}
		got := messages[1]
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Expected to have %v, got %v", want, messages[1])
		}
	})
}
