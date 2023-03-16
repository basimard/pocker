package handlers

import (
	"encoding/json"
	"net/http"
	"toggl/app/dtos"

	"github.com/sirupsen/logrus"
)

// Create a new deck
func (d *DeckHandlerImpl) CreateNewDeckHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	shuffle := query.Get("shuffle") == "true"
	cards := query.Get("cards")
	deck, err := d.deckservice.CreateNewDeck(shuffle, cards)
	if err != nil {
		d.logger.WithError(err).Error("Error creating new deck")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	writeCreateNewDeckHandlerResponse(w, deck, d.logger)
}

func writeCreateNewDeckHandlerResponse(w http.ResponseWriter, deck *dtos.RespCreateDeck, logger *logrus.Logger) {
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the deck object into JSON
	jsonData, err := json.Marshal(deck)
	if err != nil {
		logger.WithError(err).Error("Error marshaling JSON response")
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the HTTP writer
	_, err = w.Write(jsonData)
	if err != nil {
		logger.WithError(err).Error("Error writing response")
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
