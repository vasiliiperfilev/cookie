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
	router.Handle("/v1/users", http.HandlerFunc(a.userHandler))
	router.Handle("/v1/tokens", http.HandlerFunc(a.tokenHandler))

	a.Handler = a.setAccessControlHeaders(router)
	return a
}

func (a *Application) GetState() data.AppState {
	return data.AppState{Status: "available", Env: a.config.Env, Version: 1}
}
