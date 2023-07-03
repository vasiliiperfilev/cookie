package data_test

import (
	"fmt"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestPermissionModelIntegration(t *testing.T) {
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
	model := data.NewPsqlPermissionModel(db)
	t.Run("it gets correct supplier permissions", func(t *testing.T) {
		got, err := model.GetAllForType(data.UserTypeSupplier)
		tester.AssertNoError(t, err)
		want := data.Permissions{
			data.PermissionAcceptOrder,
			data.PermissionDeclineOrder,
			data.PermissionFulfillOrder,
			data.PermissionSupplierChangesOrder,
		}
		data.AssertUserPermissions(t, got, want)
	})

	t.Run("it gets correct client permissions", func(t *testing.T) {
		got, err := model.GetAllForType(data.UserTypeClient)
		tester.AssertNoError(t, err)
		want := data.Permissions{
			data.PermissionCreateOrder,
			data.PermissionClientChangesOrder,
			data.PermissionConfirmFulfillOrder,
		}
		data.AssertUserPermissions(t, got, want)
	})
}
