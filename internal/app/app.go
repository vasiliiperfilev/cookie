package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

const JsonContentType = "application/json"

type Config struct {
	Port int
	Env  string
	Db   DbConfig
}

type DbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Application struct {
	config Config
	logger *log.Logger
	models data.Models
	http.Handler
}

func New(config Config, logger *log.Logger, models data.Models) *Application {
	a := new(Application)
	a.config = config
	a.logger = logger
	a.models = models

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(a.notFoundResponse))
	router.Handle("/v1/healthcheck", http.HandlerFunc(a.healthcheckHandler))
	router.Handle("/v1/auth/register", http.HandlerFunc(a.authRegisterHandler))

	a.Handler = router
	return a
}

func (a *Application) GetState() data.AppState {
	return data.AppState{Status: "available", Env: a.config.Env, Version: 1}
}

func OpenDB(dbCfg DbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbCfg.Dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(dbCfg.MaxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(dbCfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ping to test connection
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
