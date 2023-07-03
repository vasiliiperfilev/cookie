package app_test

import (
	"log"
	"os"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

func TestOrderPost(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it POST order with correct values", func(t *testing.T) {

	})

	t.Run("it 422 if POST order with incorrect values", func(t *testing.T) {

	})

	t.Run("it 401 if POST order unathorized", func(t *testing.T) {

	})
}

func TestOrderGet(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it GET order", func(t *testing.T) {

	})

	t.Run("it 404 if GET non-existing order", func(t *testing.T) {

	})

	t.Run("it 401 if GET order unathorized", func(t *testing.T) {

	})

	t.Run("it 403 if GET not own order", func(t *testing.T) {

	})
}

func TestOrderGetAll(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it GET all orders of own id", func(t *testing.T) {

	})

	t.Run("it 404 if GET all orders of non-existing", func(t *testing.T) {

	})

	t.Run("it 401 if GET all orders uathorized", func(t *testing.T) {

	})

	t.Run("it 403 if GET all orders of not owning user", func(t *testing.T) {

	})
}

func TestOrderPut(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	// Order states with owner:
	// Created - 1 (client)
	// Accepted - 2 (supplier)
	// Declined - 3 (supplier)
	// Fulfilled - 4 (supplier)
	// Confirmed fulfillment - 5 (client)
	// Supplier changes - 6 (supplier)
	// Client changes - 7 (client)

	t.Run("it 201 if supplier PUT order to accepted", func(t *testing.T) {

	})

	t.Run("it 201 if supplier PUT order to fulfilled", func(t *testing.T) {

	})

	t.Run("it 201 if client PUT order to confirm fulfielment", func(t *testing.T) {

	})

	t.Run("it 201 if supplier PUT order to suggest changes before fulfielment", func(t *testing.T) {

	})

	t.Run("it 201 if client PUT order to suggest changes before fulfielment", func(t *testing.T) {

	})

	t.Run("it 201 if client PUT order to accept supplier changes", func(t *testing.T) {

	})

	t.Run("it 201 if client PUT order to accept client changes", func(t *testing.T) {

	})

	t.Run("it 201 if supplier PUT order to decline before accepting/suggesting changes", func(t *testing.T) {

	})

	t.Run("it 400 if supplier PUT order to decline after accepted", func(t *testing.T) {

	})

	t.Run("it 400 if supplier PUT order to approve own changes", func(t *testing.T) {

	})

	t.Run("it 400 if client PUT order to approve own changes", func(t *testing.T) {

	})

	t.Run("it 400 if supplier PUT order to canceled by client", func(t *testing.T) {

	})

	t.Run("it 400 if client PUT order to accepted by supplier", func(t *testing.T) {

	})

	t.Run("it 400 if client PUT order to declined by supplier", func(t *testing.T) {

	})

	t.Run("it 404 if PUT order of non-existing id", func(t *testing.T) {

	})

	t.Run("it 401 if PUT order unathorized", func(t *testing.T) {

	})

	t.Run("it 403 if PUT order of not owning user", func(t *testing.T) {

	})
}
