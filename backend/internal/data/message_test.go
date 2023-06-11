package data_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestMessageModelIntegration(t *testing.T) {
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
	// conversationModel := data.NewPsqlConversationModel(db)
	messageModel := data.NewPsqlMessageModel(db)
	t.Run("it doesn't insert a message if no conversation", func(t *testing.T) {
		msg := data.Message{
			ConversationId: 10,
			SenderId:       1,
			Content:        "test",
			PrevMessageId:  0,
		}
		err := messageModel.Insert(&msg)
		tester.AssertError(t, err)
		tester.AssertValue(t, err.Error(), `pq: insert or update on table "messages" violates foreign key constraint "fk_conversation_id"`, "expected no users error")
	})

	t.Run("it doesn't insert a message if user is not in conversation", func(t *testing.T) {
		// TODO: Add the test
	})

	t.Run("it inserts a message if conversation exist", func(t *testing.T) {
		msg := data.Message{
			ConversationId: 0,
			SenderId:       1,
			Content:        "test",
			PrevMessageId:  0,
		}
		err := messageModel.Insert(&msg)
		tester.AssertNoError(t, err)
	})

	t.Run("it gets messages for conversation id", func(t *testing.T) {
		want := data.Message{
			ConversationId: 0,
			SenderId:       1,
			Content:        "test get",
			PrevMessageId:  0,
		}
		err := messageModel.Insert(&want)
		tester.AssertNoError(t, err)
		messages, err := messageModel.GetAllById(int64(0))
		tester.AssertNoError(t, err)
		found := false
		for _, message := range messages {
			if reflect.DeepEqual(want, message) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to have %v message, but not found", want)
		}
	})
}
