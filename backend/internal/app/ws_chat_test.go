package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
	"golang.org/x/exp/slices"
)

type PostMessageEvent struct {
	Type    string
	Payload data.PostMessageDto
}

type MessageEvent struct {
	Type    string
	Payload data.Message
}

func TestChat(t *testing.T) {
	t.Run("2 users: sends-receive, send-receive", func(t *testing.T) {
		numUsers := 2
		cfg := app.Config{Port: 4000, Env: "development"}
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		conversationModel := data.NewStubConversationModel(generateConversation(numUsers))
		messageModel := data.NewStubMessageModel(conversationModel, []data.Message{})
		userModel := data.NewStubUserModel(generateUsers(numUsers))
		models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
		appServer := app.New(cfg, logger, models)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("2", 26))
		defer ws2.Close()
		// send first message
		conversationId := int64(1)
		dto := PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: conversationId,
				Content:        "test1",
				PrevMessageId:  0,
			},
		}
		want := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: conversationId,
			Content:        "test1",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, dto)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, conversationId, want)
		// receive first message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want) })
		// TODO: send ws read event
		assertConversationLastMessage(t, messageModel, conversationId, want.Id)
		assertUserLastReadMessage(t, conversationModel, 2, conversationId, want)
		// send second message
		dto = PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: 1,
				Content:        "test2",
				PrevMessageId:  1,
			},
		}
		want = data.Message{
			Id:             2,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test2",
			PrevMessageId:  1,
		}
		js = createWsPayload(t, dto)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, conversationId, want)
		// receive second message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want) })
		// TODO: send ws read event
		assertConversationLastMessage(t, messageModel, conversationId, want.Id)
		assertUserLastReadMessage(t, conversationModel, 2, conversationId, want)
	})

	t.Run("2 users: sends-receive, receive-send", func(t *testing.T) {
		numUsers := 2
		cfg := app.Config{Port: 4000, Env: "development"}
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		conversationModel := data.NewStubConversationModel(generateConversation(numUsers))
		messageModel := data.NewStubMessageModel(conversationModel, []data.Message{})
		userModel := data.NewStubUserModel(generateUsers(numUsers))
		models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
		appServer := app.New(cfg, logger, models)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("2", 26))
		defer ws2.Close()
		// send message from 1 to 2
		conversationId := int64(1)
		dto1 := PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: conversationId,
				Content:        "test3",
				PrevMessageId:  0,
			},
		}
		want1 := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: conversationId,
			Content:        "test3",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, dto1)
		writeWSMessage(t, ws1, js)
		assertContainsMessage(t, messageModel, conversationId, want1)
		assertConversationLastMessage(t, messageModel, conversationId, want1.Id)
		// send message from 2 to 1
		dto2 := PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: 1,
				Content:        "test4",
				PrevMessageId:  1,
			},
		}
		want2 := data.Message{
			Id:             2,
			SenderId:       2,
			ConversationId: 1,
			Content:        "test4",
			PrevMessageId:  1,
		}
		js = createWsPayload(t, dto2)
		writeWSMessage(t, ws2, js)
		assertContainsMessage(t, messageModel, 1, want2)
		assertConversationLastMessage(t, messageModel, conversationId, want2.Id)
		// user 2: receive message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws2, want1) })
		// TODO: send read message event
		// user 1: receive message
		within(t, 500*time.Millisecond, func() { assertMessage(t, ws1, want2) })
		// TODO: send read message event
		assertUserLastReadMessage(t, conversationModel, 1, conversationId, want2)
		assertUserLastReadMessage(t, conversationModel, 2, conversationId, want2)
	})

	t.Run("9 users send messages to 1 user: concurrent sends", func(t *testing.T) {
		messageModel, appServer := createServer(10)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
		// send messages from 1 to 2
		for i := 2; i <= 9; i++ {
			ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat(strconv.Itoa(i), 26))
			defer ws2.Close()
			for j := 1; j <= 100; j++ {
				dto := PostMessageEvent{
					Type: app.EventMessage,
					Payload: data.PostMessageDto{
						ConversationId: int64(i - 1),
						Content:        fmt.Sprintf("test%v", i),
						PrevMessageId:  0,
					},
				}
				want := data.Message{
					Id:             int64(j),
					SenderId:       int64(i),
					ConversationId: int64(i - 1),
					Content:        fmt.Sprintf("test%v", i),
					PrevMessageId:  0,
				}
				js := createWsPayload(t, dto)
				writeWSMessage(t, ws2, js)
				assertContainsMessage(t, messageModel, int64(i-1), want)
				within(t, 500*time.Millisecond, func() { assertMessage(t, ws1, want) })
			}
		}
	})

	t.Run("only user in conversation receive messages", func(t *testing.T) {
		cfg := app.Config{Port: 4000, Env: "development"}
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		// create 3 users
		userModel := data.NewStubUserModel(generateUsers(3))
		// create only 1 conversation between user ids 1 and 2
		conversationModel := data.NewStubConversationModel(generateConversation(2))
		messageModel := data.NewStubMessageModel(conversationModel, []data.Message{})
		models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
		appServer := app.New(cfg, logger, models)
		server := httptest.NewServer(appServer)
		conns := []*websocket.Conn{}
		for i := 1; i <= 3; i++ {
			ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat(strconv.Itoa(i), 26))
			defer ws.Close()
			conns = append(conns, ws)
		}
		// send a message
		dto := PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: 1,
				Content:        "test1",
				PrevMessageId:  0,
			},
		}
		want := data.Message{
			Id:             1,
			SenderId:       1,
			ConversationId: 1,
			Content:        "test1",
			PrevMessageId:  0,
		}
		js := createWsPayload(t, dto)
		writeWSMessage(t, conns[0], js)
		assertContainsMessage(t, messageModel, 1, want)
		// user 2 receives message
		within(t, 500*time.Millisecond, func() { assertMessage(t, conns[1], want) })
		// sender receive a message back as confirmation
		within(t, 500*time.Millisecond, func() { assertMessage(t, conns[0], want) })
		// user 3 doesn't receive the message
		assertNoMessage(t, conns[2])
	})
}

func TestChatErrors(t *testing.T) {
	t.Run("it handles client disconnection", func(t *testing.T) {
		_, appServer := createServer(2)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		err := ws1.Close()
		tester.AssertNoError(t, err)
		// can establish ws connection again
		ws1 = mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
	})

	t.Run("it handles client disconnection during processing", func(t *testing.T) {
		_, appServer := createServer(2)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("2", 26))
		// send first message
		want := PostMessageEvent{
			Type: app.EventMessage,
			Payload: data.PostMessageDto{
				ConversationId: 1,
				Content:        "test1",
				PrevMessageId:  0,
			},
		}
		js := createWsPayload(t, want)
		writeWSMessage(t, ws1, js)
		err := ws2.Close()
		tester.AssertNoError(t, err)
		err = ws1.Close()
		tester.AssertNoError(t, err)
		// can establish ws connection gain
		ws1 = mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
		ws2 = mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("2", 26))
		defer ws2.Close()
	})

	t.Run("it responds with error event if payload is incorrect", func(t *testing.T) {
		// client1 connects
		// client2 connects
		// client1 sends incorrect message to client 2
		// client1 receives incorrect event
		// client 2 receives nothing
		_, appServer := createServer(2)
		server := httptest.NewServer(appServer)
		defer server.Close()
		ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("1", 26))
		defer ws1.Close()
		ws2 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat?token="+strings.Repeat("2", 26))
		defer ws2.Close()
		// send first message
		want := struct {
			Type    string
			Payload string
		}{
			Type:    app.EventMessage,
			Payload: "wrong payload",
		}
		js := createWsPayload(t, want)
		writeWSMessage(t, ws1, js)
		// receive error back message
		wantError := app.ErrorResponse{Message: app.PayloadErrorMessage, Errors: map[string]string{}}
		within(t, 500*time.Millisecond, func() { assertErrorEvent(t, ws1, wantError) })
		// client 2 receives nothing
		assertNoMessage(t, ws2)
	})

	// TODO: uncomment and finish up after WS is extracted as separate package
	// t.Run("it closes connection if no pong response", func(t *testing.T) {
	// 	_, appServer := createServer(2)
	// 	server := httptest.NewServer(appServer)
	// 	defer server.Close()
	// 	h1 := http.Header{"Authorization": {"Bearer " + strings.Repeat("1", 26)}}
	// 	ws1 := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/v1/chat", h1)
	// 	ws1.SetPingHandler(func(appData string) error { return nil })
	// 	time.Sleep(10)
	// 	_, _, err := ws1.ReadMessage()
	// 	tester.AssertError(t, err)
	// })
}

func createServer(numUsers int) (*data.StubMessageModel, *app.Application) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	conversationModel := data.NewStubConversationModel(generateConversation(numUsers))
	messageModel := data.NewStubMessageModel(conversationModel, []data.Message{})
	userModel := data.NewStubUserModel(generateUsers(numUsers))
	models := data.Models{Message: messageModel, User: userModel, Conversation: conversationModel}
	appServer := app.New(cfg, logger, models)
	return messageModel, appServer
}

func generateConversation(numUsers int) []data.Conversation {
	c := []data.Conversation{}
	id := 1
	for i := 1; i <= numUsers; i++ {
		for j := i + 1; j <= numUsers; j++ {
			c = append(c, data.Conversation{Id: int64(id), UserIds: []int64{int64(i), int64(j)}, LastReadMessages: []data.LastReadMessage{}})
			id++
		}
	}
	return c
}

func generateUsers(numUsers int) []data.User {
	u := []data.User{}
	for i := 1; i <= numUsers; i++ {
		u = append(u, data.User{Id: int64(i), Name: fmt.Sprintf("test user %v", i), Type: i%2 + 1})
	}
	return u
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

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

func assertContainsMessage(t *testing.T, m data.MessageModel, conversationId int64, want data.Message) {
	t.Helper()

	passed := tester.RetryUntil(500*time.Millisecond, func() bool {
		messages, err := m.GetAllByConversationId(int64(conversationId))
		tester.AssertNoError(t, err)
		return slices.Contains(messages, want)
	})

	if !passed {
		t.Fatalf("Expected to have %v", want)
	}
}

func assertMessage(t *testing.T, ws *websocket.Conn, want data.Message) {
	t.Helper()

	passed := tester.RetryUntil(1000*time.Millisecond, func() bool {
		_, msg, err := ws.ReadMessage()
		tester.AssertNoError(t, err)
		var got MessageEvent
		json.NewDecoder(bytes.NewReader(msg)).Decode(&got)
		return reflect.DeepEqual(got.Payload, want)
	})

	if !passed {
		t.Fatalf("Expected to have %v", want)
	}
}

func assertNoMessage(t *testing.T, ws *websocket.Conn) {
	t.Helper()

	done := make(chan []byte, 1)

	go func() {
		_, msg, _ := ws.ReadMessage()
		done <- msg
	}()

	select {
	case msg := <-done:
		t.Errorf("Got message %s, expected nothing", string(msg))
	case <-time.After(500 * time.Millisecond):
	}
}

func assertErrorEvent(t *testing.T, ws *websocket.Conn, want app.ErrorResponse) {
	t.Helper()

	_, msg, err := ws.ReadMessage()
	tester.AssertNoError(t, err)
	var gotEvent app.WsEvent
	json.NewDecoder(bytes.NewReader(msg)).Decode(&gotEvent)
	var gotPayload app.ErrorResponse
	json.NewDecoder(bytes.NewReader(gotEvent.Payload)).Decode(&gotPayload)

	passed := tester.RetryUntil(1000*time.Millisecond, func() bool {
		return reflect.DeepEqual(gotEvent.Type, app.EventError) && reflect.DeepEqual(gotPayload, want)
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

func assertUserLastReadMessage(t *testing.T, conversationModel *data.StubConversationModel, userId int64, conversationId int64, want data.Message) {
	t.Helper()
	conversation, err := conversationModel.GetById(conversationId)
	tester.AssertNoError(t, err)
	exists := false
	for _, readMessage := range conversation.LastReadMessages {
		if readMessage.UserId == userId {
			exists = true
			tester.AssertValue(t, readMessage.Message, want, "Expected same last read message")
		}
	}
	if !exists {
		t.Fatalf("Last read message for user id %v not found!", userId)
	}
}

func assertConversationLastMessage(t *testing.T, messageModel *data.StubMessageModel, conversationId int64, want int64) {
	t.Helper()
	messages, err := messageModel.GetAllByConversationId(conversationId)
	tester.AssertNoError(t, err)
	got := messages[len(messages)-1].Id
	if got != want {
		t.Fatalf("Expected last read message id %v, got id %v", want, got)
	}
}
