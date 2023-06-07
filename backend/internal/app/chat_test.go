package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
	"golang.org/x/exp/slices"
)

func TestChat(t *testing.T) {
	env := "testing"
	cfg := app.Config{Port: 4000, Env: env}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	conversationModel := data.NewStubConversationModel([]data.Conversation{{Id: 1, UserIds: []int64{1, 2}}})
	messageModel := data.NewStubMessageModel([]data.Conversation{{Id: 1, UserIds: []int64{1, 2}}}, []data.Message{})
	userModel := data.NewStubUserModel([]data.User{{Id: 1}, {Id: 2}})
	models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
	appServer := app.New(cfg, logger, models)
	t.Run("sends and receives a message", func(t *testing.T) {
		server := httptest.NewServer(appServer)
		defer server.Close()
		h1 := http.Header{"Authorization": {"Bearer " + strings.Repeat("1", 26)}}
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat", h1)
		defer ws1.Close()
		h2 := http.Header{"Authorization": {"Bearer " + strings.Repeat("2", 26)}}
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat", h2)
		defer ws2.Close()
		// send first message
		want := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, want)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, want)
		// receive first message
		within(t, 500*time.Microsecond, func() { assertMessage(t, ws2, want) })
		// send second message
		want = data.Message{
			Id:             2,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test2",
			PrevMessageId:  1,
		}
		js = createWsPayload(t, want)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, want)
		// receive second message
		within(t, 500*time.Microsecond, func() { assertMessage(t, ws2, want) })
	})
}

func mustDialWS(t *testing.T, url string, headers http.Header) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, headers)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message []byte) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func createWsPayload(t *testing.T, payload any) []byte {
	js, err := json.Marshal(payload)
	tester.AssertNoError(t, err)
	return js
}

func assertContainsMessage(t *testing.T, m data.MessageModel, want data.Message) {
	t.Helper()

	passed := tester.RetryUntil(500*time.Millisecond, func() bool {
		messages, err := m.GetAllByUserId(1)
		tester.AssertNoError(t, err)
		return slices.Contains(messages, want)
	})

	if !passed {
		t.Fatalf("Expected to have %v", want)
	}
}

func assertMessage(t *testing.T, ws *websocket.Conn, want data.Message) {
	t.Helper()

	_, msg, err := ws.ReadMessage()
	tester.AssertNoError(t, err)
	var got data.Message
	json.NewDecoder(bytes.NewReader(msg)).Decode(&got)

	passed := tester.RetryUntil(500*time.Millisecond, func() bool {
		return reflect.DeepEqual(got, want)
	})

	if !passed {
		t.Fatalf("Expected to have %v", want)
	}
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
