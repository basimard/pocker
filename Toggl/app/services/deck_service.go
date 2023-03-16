package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"toggl/app/dtos"
	"toggl/app/models"
	"toggl/app/repos"

	"github.com/sirupsen/logrus"
)

// suits and values of cards type
var suits = []string{"SPADES", "HEARTS", "DIAMONDS", "CLUBS"}
var values = []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}

type DeckService interface {
	CreateNewDeck(shuffled bool, cards string) (*dtos.RespCreateDeck, error)
	OpenDeck(deckId string) (*dtos.RespOpenDeck, error)
	DrawCard(deckId string, count int) (*dtos.RespDrawDeck, error)
}

type DeckServiceImpl struct {
	logger *logrus.Logger
	repo   *repos.Repository
}

// New Deck service setup using dependencies
func NewDeckService(logger *logrus.Logger, repo *repos.Repository) *DeckServiceImpl {
	return &DeckServiceImpl{logger: logger, repo: repo}
}

// parse cards and validate for creating deck
func parseCode(code string) (models.Card, error) {
	if len(code) != 2 {
		return models.Card{}, fmt.Errorf("%s is not a valid code: must be exactly two letters", code)
	}

	value := ""
	for _, v := range values {
		if v[0] == code[0] {
			value = v
			break
		}
	}

	if value == "" {
		return models.Card{}, fmt.Errorf("%s is not a valid code: unknown value %s", code, code[0])
	}

	suit := ""
	for _, s := range suits {
		if s[0] == code[1] {
			suit = s
			break
		}
	}

	if suit == "" {
		return models.Card{}, fmt.Errorf("%s is not a valid code: unknown suit %s", code, code[1])
	}

	return models.Card{Value: value, Suit: suit, Code: code}, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Create a full deck with 52 cards
func CreateFullDeck() []models.Card {

	var deckCards []models.Card

	for _, suit := range suits {
		for _, value := range values {
			var code = strings.ToUpper(suit[:1]) + strings.ToUpper(value[:1])
			deckCards = append(deckCards, models.Card{Suit: suit, Value: value, Code: code})
		}
	}

	return deckCards
}

// create a new deck with params
func (s *DeckServiceImpl) CreateNewDeck(shuffled bool, cards string) (*dtos.RespCreateDeck, error) {

	var lstCards = strings.Split(cards, ",")
	var deckCards []models.Card
	if len(lstCards) != 1 && lstCards[0] != "" {

		if len(lstCards) > 52 {
			s.logger.Errorf("Number of cards exceeded")
		}
		for _, code := range lstCards {
			parsedCard, err := parseCode(code)
			if err != nil {
				s.logger.Errorf("%s is not a valid code\n", code)
				continue
			}

			deckCards = append(deckCards, parsedCard)
		}

	} else {
		deckCards = CreateFullDeck()
	}

	rand.Seed(time.Now().UnixNano())
	if shuffled == true {
		rand.Shuffle(len(deckCards), func(i, j int) {
			deckCards[i], deckCards[j] = deckCards[j], deckCards[i]
		})
	}

	deck := &models.Deck{
		Shuffled:  shuffled,
		Remaining: len(deckCards),
		Cards:     deckCards,
	}

	var result, _ = s.repo.CreateDeck(deck)

	var resp = dtos.RespCreateDeck{DeckID: result, Remaining: deck.Remaining, Shuffled: deck.Shuffled}

	return &resp, nil
}

// open a new deck based on id
func (s *DeckServiceImpl) OpenDeck(deckId string) (*dtos.RespOpenDeck, error) {
	exist, err := s.repo.CheckDeckExist(deckId)
	if err != nil {
		s.logger.Errorf("Error in checking id %s", deckId)
		return nil, err
	}
	if !exist {
		s.logger.Errorf("Deck with id %s does not exist", deckId)
		return nil, err
	}

	deck, err := s.repo.OpenDeck(deckId)
	if err != nil {
		return nil, err
	}

	return deck, nil
}

// Draw number of cards from deck based on id
func (s *DeckServiceImpl) DrawCard(deckId string, count int) (*dtos.RespDrawDeck, error) {
	// Check if deck exists
	exist, err := s.repo.CheckDeckExist(deckId)
	if err != nil {
		s.logger.Errorf("Error in checking id %s", deckId)
		return nil, err
	}
	if !exist {
		s.logger.Errorf("Deck with id %s does not exist", deckId)
		return nil, err
	}

	// Check if count is less than remaining cards
	deck, err := s.repo.OpenDeck(deckId)
	if err != nil {
		return nil, err
	}
	remaining := len(deck.Cards)
	if count > remaining {
		s.logger.Errorf("Requested count %d exceeds remaining cards %d in deck", count, remaining)
		return nil, err
	}

	// Draw cards
	cards, err := s.repo.DrawCard(deckId, count)
	if err != nil {
		s.logger.Errorf("Error in draw %d cards from deck %s", count, remaining)
		return nil, err
	}

	return cards, nil
}
