package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type AppConfig struct {
	ServerHost string	`envconfig:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort int 		`envconfig:"SERVER_PORT" default:"8080"`
}

func main() {
	logger := zerolog.New(zerolog.MultiLevelWriter(os.Stdout)).
		With().
		Timestamp().
		Logger()

	var config AppConfig
	err := envconfig.Process("myapp", &config)
	if err != nil {
		logger.Fatal().Err(err)
	}

	userConfig := user.NewConfig()
	userClient := user.NewClient(userConfig)
	h := NewHandler(userClient, logger)

	r := chi.NewRouter()
	r.Use(h.MiddlewareLogger)
	r.Use(middleware.Recoverer)
	r.Get("/v1/user-posts/{id}", h.GetUserPostsHandler)

	s := http.Server {
		Addr: fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		Handler: r,
	}
	logger.Info().Msgf("server listening on port %d", config.ServerPort)
	err = s.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err)
	}
}

