package dtos

type RespDrawDeck struct {
	Cards []RespDrawCard `json:"cards"`
}

type RespDrawCard struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}
