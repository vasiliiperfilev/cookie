package data_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestOrderModelIntegration(t *testing.T) {
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
	t.Run("it inserts Order and retrieves it", func(t *testing.T) {
		orderModel := data.NewPsqlOrderModel(db)
		dto := data.PostOrderDto{
			ConversationId: 1,
			ClientId:       2,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 3,
				},
			},
		}
		want, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		got, err := orderModel.GetById(want.Id)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, got, want, "Expected same item")
	})

	t.Run("it retrieves orders by user id", func(t *testing.T) {
		orderModel := data.NewPsqlOrderModel(db)
		dto := data.PostOrderDto{
			ConversationId: 2,
			ClientId:       4,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 3,
				},
			},
		}
		// insert 2 orders
		want1, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		want2, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		want := []data.Order{want1, want2}
		got, err := orderModel.GetAllByUserId(4)
		tester.AssertNoError(t, err)
		for i := range got {
			tester.AssertValue(t, got[i], want[i], "Expected same item")
		}
	})

	t.Run("it updates order", func(t *testing.T) {
		orderModel := data.NewPsqlOrderModel(db)
		dto := data.PostOrderDto{
			ConversationId: 1,
			ClientId:       2,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 3,
				},
			},
		}
		order, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		order.Items = []data.ItemQuantity{
			{
				ItemId:   2,
				Quantity: 3,
			},
		}
		order.StateId = data.OrderStateClientChanges
		time.Sleep(1 * time.Second)
		want, err := orderModel.Update(order)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, want, order, "Expected same item from update order")
		got, err := orderModel.GetById(want.Id)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, got, want, "Expected same item from get order")
	})
}
