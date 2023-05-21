package app

import (
	"log"
	"net/http"
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

type State struct {
	Status  string
	Env     string
	Version int
}

func New(config Config, logger *log.Logger) *Application {
	a := new(Application)
	a.config = config
	a.logger = logger

	router := http.NewServeMux()
	router.Handle("/v1/healthcheck", http.HandlerFunc(a.healthcheckHandler))

	a.Handler = router
	return a
}

func (a *Application) GetState() State {
	return State{Status: "available", Env: a.config.Env, Version: 1}
}
