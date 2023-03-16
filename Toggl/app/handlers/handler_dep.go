package handlers

import (
	"toggl/app/services"

	"github.com/sirupsen/logrus"
)

type DeckHandlerImpl struct {
	deckservice services.DeckService
	logger      *logrus.Logger
}

// Setup a new DeckHandler with deck service and logger
func NewDeckHandler(deckService services.DeckService, logger *logrus.Logger) *DeckHandlerImpl {
	return &DeckHandlerImpl{deckservice: deckService, logger: logger}
}
