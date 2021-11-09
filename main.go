package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hooliganlin/simple-go-rest-api/cache"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

type AppConfig struct {
	ServerHost			string			`envconfig:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort			int				`envconfig:"SERVER_PORT" default:"8080"`
	CacheTTL			time.Duration	`envconfig:"CACHE_TTL" default:"5m"`
	CacheTTLInterval	time.Duration	`envconfig:"CACHE_TTL_INTERVAL" default:"10m"`
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

	//set cache
	c := cache.NewDefaultCache(config.CacheTTL, config.CacheTTLInterval)

	userConfig := user.NewConfig()
	userClient := user.NewDefaultClient(userConfig, c)
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
	if err = s.ListenAndServe(); err != nil {
		logger.Fatal().Err(err)
	}
}

