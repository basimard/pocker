package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"toggl/app/dtos"
	"toggl/app/utils"

	"github.com/sirupsen/logrus"
)

// Open a deck handler
func (d *DeckHandlerImpl) OpenDeckHandler(w http.ResponseWriter, r *http.Request) {

	// Get the deck ID from the URL parameter

	deckId := r.URL.Query().Get("deck_id")
	fmt.Println(deckId)
	if deckId == "" {
		d.logger.Error("Empty deck id")
		http.Error(w, fmt.Sprintf("Deck id parameter is required"), http.StatusBadRequest)
		return
	}

	_, err := utils.Parse_uuid(deckId)
	if err != nil {
		d.logger.WithError(err).Error("Error in parsing deck id ")
		http.Error(w, fmt.Sprintf("Invalid deck id"), http.StatusBadRequest)
		return
	}
	// Fetch the deck by its ID
	deck, err := d.deckservice.OpenDeck(deckId)
	if err != nil {
		d.logger.WithError(err).Error("Error in open deck ")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	writeOpenDeckHandlerResponse(w, deck, d.logger)
}

func writeOpenDeckHandlerResponse(w http.ResponseWriter, deck *dtos.RespOpenDeck, logger *logrus.Logger) {
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
