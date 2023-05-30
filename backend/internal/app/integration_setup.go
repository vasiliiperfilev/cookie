package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

func PrepareServer(db *sql.DB, port int) *Application {
	cfg := Config{Port: port, Env: "development"}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.NewModels(db)

	server := New(cfg, logger, models)
	return server
}
