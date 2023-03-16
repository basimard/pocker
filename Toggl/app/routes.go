package app

import (
	"toggl/app/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(mux *mux.Router, deckHandler *handlers.DeckHandlerImpl) {
	// Register the handlers with the HTTP server
	mux.HandleFunc("/v1/create-deck", deckHandler.CreateNewDeckHandler).Methods("POST")
	mux.HandleFunc("/v1/open-deck", deckHandler.OpenDeckHandler).Methods("GET")
	mux.HandleFunc("/v1/draw-cards", deckHandler.DrawCardHandler).Methods("POST")
}
