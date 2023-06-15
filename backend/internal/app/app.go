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
	hub    *Hub
	http.Handler
}

func New(config Config, logger *log.Logger, models data.Models) *Application {
	a := new(Application)
	a.config = config
	a.logger = logger
	a.models = models
	// start websocket hub
	a.hub = newHub(a)
	go a.hub.run()
	// create router
	router := a.routes()
	a.Handler = a.setAccessControlHeaders(router)

	return a
}

func (a *Application) GetState() data.AppState {
	return data.AppState{Status: "available", Env: a.config.Env, Version: 1}
}
