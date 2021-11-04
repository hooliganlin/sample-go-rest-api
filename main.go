package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type AppConfig struct {
	ServerHost string	`envconfig:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort int 		`envconfig:"SERVER_PORT" default:"8080"`
}

func main() {
	var config AppConfig
	err := envconfig.Process("myapp", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	userConfig := user.NewConfig()
	userClient := user.NewClient(userConfig)
	h := NewHandler(userClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/v1/user-posts/{id}", h.GetUserPostsHandler)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort), r)
	if err != nil {
		log.Fatal(err.Error())
	}
}