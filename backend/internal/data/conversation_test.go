package data_test

import (
	"fmt"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
	"golang.org/x/exp/slices"
)

func TestConversationModelIntegration(t *testing.T) {
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
	t.Run("it doesn't insert a conversation if no users", func(t *testing.T) {
		model := data.NewPsqlConversationModel(db)
		dto := data.PostConversationDto{
			UserIds: []int64{99, 100},
		}
		_, err := model.Insert(dto)
		tester.AssertError(t, err)
		tester.AssertValue(t, err.Error(), `pq: insert or update on table "conversations_users" violates foreign key constraint "conversations_users_user_id_fkey"`, "expected no users error")
	})

	t.Run("it inserts a conversation if users exist", func(t *testing.T) {
		model := data.NewPsqlConversationModel(db)
		dto := data.PostConversationDto{
			UserIds: []int64{1, 2},
		}
		_, err := model.Insert(dto)
		tester.AssertNoError(t, err)
	})

	t.Run("it gets conversations list for user id", func(t *testing.T) {
		model := data.NewPsqlConversationModel(db)
		userId := int64(1)
		conversations, err := model.GetAllByUserId(userId)
		tester.AssertNoError(t, err)
		for _, conversation := range conversations {
			userIds := data.Map(conversation.Users, func(u data.User) int64 { return u.Id })
			if !slices.Contains(userIds, userId) {
				t.Fatalf("Expected to have user id %v in conversation", userId)
			}
		}
	})
}
