package models

type Deck struct {
	DeckID    string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}
