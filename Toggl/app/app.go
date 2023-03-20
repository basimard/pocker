package app

import (
	"fmt"
	"log"

	"net/http"
	"toggl/app/config"
	"toggl/app/handlers"
	"toggl/app/repos"
	"toggl/app/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type App struct {
	httpServer *http.Server
}

func NewApp(config *config.Config) (*App, error) {

	// Create a new HTTP server with the desired configuration
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	logger := logrus.New()

	deckRepo := repos.NewRepository(logger, false, config)
	// Create new services for the app
	deckService := services.NewDeckService(logger, deckRepo)

	// Create new handlers for the app, injecting the services
	deckHandler := handlers.NewDeckHandler(deckService, logger)

	// Create a new ServeMux object
	mux := mux.NewRouter()

	// Register the routes with the ServeMux object
	RegisterRoutes(mux, deckHandler)

	// Attach the ServeMux to the HTTP server
	httpServer.Handler = mux

	return &App{httpServer: httpServer}, nil

}

func (a *App) Start() error {
	log.Printf("Starting server on %s", a.httpServer.Addr)

	// Start the HTTP server
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *App) Stop() error {
	log.Printf("Stopping server on %s", a.httpServer.Addr)

	// Shutdown the HTTP server gracefully
	err := a.httpServer.Shutdown(nil)
	if err != nil {
		return err
	}

	return nil
}
