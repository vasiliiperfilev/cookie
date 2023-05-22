package app

import (
	"log"
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

const JsonContentType = "application/json"

type Config struct {
	Port int
	Env  string
}

type Application struct {
	config Config
	logger *log.Logger
	http.Handler
}

func New(config Config, logger *log.Logger) *Application {
	a := new(Application)
	a.config = config
	a.logger = logger

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
