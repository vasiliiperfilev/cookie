package data_test

import (
	"fmt"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestItemModelIntegration(t *testing.T) {
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
	supplierId := int64(5)
	testData := []data.Item{
		{
			SupplierId: supplierId,
			Unit:       "l",
			Size:       1,
			Name:       "Milk",
			ImageUrl:   "test",
		},
		{
			SupplierId: supplierId,
			Unit:       "kg",
			Size:       1,
			Name:       "Apples",
			ImageUrl:   "test",
		},
	}
	t.Run("it inserts Item and retrieves it", func(t *testing.T) {
		model := data.NewPsqlItemModel(db)
		err := model.Insert(&testData[0])
		tester.AssertNoError(t, err)
		got, err := model.GetById(testData[0].Id)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, got, testData[0], "Expected same item")
	})

	t.Run("it retrieves items by user id", func(t *testing.T) {
		model := data.NewPsqlItemModel(db)
		want := testData
		err := model.Insert(&want[1])
		tester.AssertNoError(t, err)
		got, err := model.GetAllBySupplierId(supplierId)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, got, want, "Expected same items array")
	})
}
