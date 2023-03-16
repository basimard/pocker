package models

type Card struct {
	Id     string `json:"id"`
	DeckId string `json:"deck_id"`
	Code   string `json:"code"`
	Value  string `json:"value"`
	Suit   string `json:"suit"`
	Drawn  int    `json:"drawn"`
}
