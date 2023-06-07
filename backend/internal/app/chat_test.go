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
	t.Run("2 users: sends-receive, send-receive", func(t *testing.T) {
		messageModel, appServer := createServer()
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
			Content:        "test1",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, want)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, 1, want)
		// receive first message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want) })
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
		assertContainsMessage(t, messageModel, 1, want)
		// receive second message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want) })
	})
	t.Run("2 users: sends-receive, receive-send", func(t *testing.T) {
		messageModel, appServer := createServer()
		server := httptest.NewServer(appServer)
		defer server.Close()
		h1 := http.Header{"Authorization": {"Bearer " + strings.Repeat("1", 26)}}
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat", h1)
		defer ws1.Close()
		h2 := http.Header{"Authorization": {"Bearer " + strings.Repeat("2", 26)}}
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat", h2)
		defer ws2.Close()
		// send message from 1 to 2
		want1 := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test3",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, want1)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, 1, want1)
		// send message from 2 to 1
		want2 := data.Message{
			Id:             2,
			SenderId:       2,
			ConversationId: 1,
			Content:        "test4",
			PrevMessageId:  1,
		}
		js = createWsPayload(t, want2)
		writeWSMessage(t, ws2, js)
		assertContainsMessage(t, messageModel, 2, want2)
		// user 2: receive message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want1) })
		// user 1: receive message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws1, want2) })
	})
}

func createServer() (*data.StubMessageModel, *app.Application) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	conversationModel := data.NewStubConversationModel([]data.Conversation{{Id: 1, UserIds: []int64{1, 2}}})
	messageModel := data.NewStubMessageModel([]data.Conversation{{Id: 1, UserIds: []int64{1, 2}}}, []data.Message{})
	userModel := data.NewStubUserModel([]data.User{{Id: 1}, {Id: 2}})
	models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
	appServer := app.New(cfg, logger, models)
	return messageModel, appServer
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

func assertContainsMessage(t *testing.T, m data.MessageModel, userId int, want data.Message) {
	t.Helper()

	passed := tester.RetryUntil(500*time.Millisecond, func() bool {
		messages, err := m.GetAllByUserId(int64(userId))
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

	passed := tester.RetryUntil(1000*time.Millisecond, func() bool {
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
