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
	t.Run("it inserts Item and retrieves it", func(t *testing.T) {
		model := data.NewPsqlItemModel(db)
		want := data.Item{
			SupplierId: 1,
			Unit:       "l",
			Size:       1,
			Name:       "Milk",
			ImageUrl:   "test",
		}
		err := model.Insert(&want)
		tester.AssertNoError(t, err)
		got, err := model.GetById(want.Id)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, got, want, "Expected same item")
	})
}
