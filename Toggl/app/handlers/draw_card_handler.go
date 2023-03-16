package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"toggl/app/dtos"
	"toggl/app/utils"

	"github.com/sirupsen/logrus"
)

// Draw a card from deck
func (d *DeckHandlerImpl) DrawCardHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	deckId := r.URL.Query().Get("deck_id")
	countStr := r.URL.Query().Get("count")

	// Validate deckId parameter
	if deckId == "" {
		d.logger.Error("Empty deck id")
		http.Error(w, "Deck id parameter is required", http.StatusBadRequest)
		return
	}
	_, err := utils.Parse_uuid((deckId))
	if err != nil {
		d.logger.WithError(err).Error("Error in parsing ")
		http.Error(w, fmt.Sprintf("Invalid deck id"), http.StatusBadRequest)
		return
	}

	// Validate count parameter
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		d.logger.WithError(err).Error("Error in checking is count positive number")
		http.Error(w, "Count parameter must be a positive integer", http.StatusBadRequest)
		return
	}

	// Call service method to draw cards
	deck, err := d.deckservice.DrawCard(deckId, count)
	if err != nil {
		d.logger.WithError(err).Error("Error in draw a card")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	writeDrawCardHandlerResponse(w, deck, d.logger)
}

func writeDrawCardHandlerResponse(w http.ResponseWriter, deck *dtos.RespDrawDeck, logger *logrus.Logger) {
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
