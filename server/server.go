// Package server provides a HTTP listener for handling requests from either a
// game server or any form of interface to the data.
package server

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/Southclaws/ScavengeSurviveCore/storage"
	"github.com/Southclaws/ScavengeSurviveCore/types"
)

// Config stores static configuration
type Config struct {
	Bind string `split_words:"true" required:"true"` // bind interface
}

// App stores and controls program state
type App struct {
	config    Config
	handlers  map[string][]Route
	validator *validator.Validate
	store     types.Storer
	ctx       context.Context
	cancel    context.CancelFunc
}

// Start fires up a HTTP server and routes API calls to the database manager
func Start(config Config) {
	logger.Debug("initialising ssc-server with debug logging", zap.Any("config", config))

	app := App{
		config:    config,
		validator: validator.New(),
		store:     storage.New(storage.Config{}),
	}
	app.handlers = map[string][]Route{
		"player": playerRoutes(app),
		"report": reportRoutes(app),
		"ban":    banRoutes(app),
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())

	router := mux.NewRouter().StrictSlash(true)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"HEAD", "GET", "POST", "PUT", "OPTIONS"})

	for name, routes := range app.handlers {
		logger.Debug("loaded handler",
			zap.String("name", name),
			zap.Int("routes", len(routes)))

		for _, route := range routes {
			if route.Authenticated {
				router.Methods(route.Method).
					Path(route.Path).
					Name(route.Name).
					Handler(app.Authenticator(route.handler))
			} else {
				router.Methods(route.Method).
					Path(route.Path).
					Name(route.Name).
					Handler(route.handler)
			}

			logger.Debug("registered handler route",
				zap.String("name", route.Name),
				zap.String("method", route.Method),
				zap.String("path", route.Path))
		}
	}

	logger.Debug("initialisation complete")
	err := http.ListenAndServe(app.config.Bind, handlers.CORS(headersOk, originsOk, methodsOk)(router))

	logger.Fatal("unexpected termination",
		zap.Error(err))
}