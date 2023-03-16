package dtos

type RespOpenDeck struct {
	DeckID    string             `json:"deck_id"`
	Shuffled  bool               `json:"shuffled"`
	Remaining int                `json:"remaining"`
	Cards     []RespOpenDeckCard `json:"cards"`
}

type RespOpenDeckCard struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}
