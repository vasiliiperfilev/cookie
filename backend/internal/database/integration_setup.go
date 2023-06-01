package database

import (
	"database/sql"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/tester"
)

const (
	POSTGRES_DB       = "cookie_testing"
	POSTGRES_USER     = "cookie"
	POSTGRES_PASSWORD = "cookie"
	POSTGRES_PORT     = "5433"
)

func PrepareTestDb(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	cfg := Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	// open connection
	db, err := OpenDB(cfg)
	tester.AssertNoError(t, err)

	return db
}
