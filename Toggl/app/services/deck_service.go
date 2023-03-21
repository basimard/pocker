package services

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
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
func parseCode(code string, logger *logrus.Logger) (*models.Card, error) {
	if len(code) != 2 {
		logger.Error("Invalid card")
		return nil, errors.New("Invalid card")
	}

	value := ""
	for _, v := range values {
		if v[0] == code[0] {
			value = v
			break
		}
	}

	if value == "" {
		logger.Error("Invalid value")
		return nil, errors.New("Invalid value")
	}

	suit := ""
	for _, s := range suits {
		if s[0] == code[1] {
			suit = s
			break
		}
	}

	if suit == "" {
		logger.Error("Invalid suit")
		return nil, errors.New("Invalid suit")
	}

	return &models.Card{Value: value, Suit: suit, Code: code}, nil
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
			s.logger.Error("Number of cards exceeded")
		}
		for _, code := range lstCards {
			parsedCard, err := parseCode(code, s.logger)
			if err != nil {
				s.logger.Error("%s is not a valid code\n", code)
				return nil, err
			}

			deckCards = append(deckCards, *parsedCard)
		}

	} else {
		deckCards = CreateFullDeck()
	}

	if shuffled == true {
		deckCards = shuffleCards(deckCards)
	}

	deck := &models.Deck{
		Shuffled:  shuffled,
		Remaining: len(deckCards),
		Cards:     deckCards,
	}

	result, err := s.repo.CreateDeck(deck)
	if err != nil {
		s.logger.WithError(err).Error("Error in creating deck")
	}

	var resp = dtos.RespCreateDeck{DeckID: result, Remaining: deck.Remaining, Shuffled: deck.Shuffled}

	return &resp, nil
}

// shuffle the cards
func shuffleCards(deck []models.Card) []models.Card {
	n := len(deck)
	for i := n - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		deck[i], deck[j.Int64()] = deck[j.Int64()], deck[i]
	}
	return deck
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
		return nil, errors.New("Id doesn't exist")
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
		return nil, errors.New("Id doesn't exist")
	}

	// Check if count is less than remaining cards
	deck, err := s.repo.OpenDeck(deckId)
	if err != nil {
		return nil, err
	}
	remaining := len(deck.Cards)
	if count > remaining {
		s.logger.Errorf("Requested count %d exceeds remaining cards %d in deck", count, remaining)
		return nil, errors.New("Requested count exceeds remaining cards in deck")
	}

	// Draw cards
	cards, err := s.repo.DrawCard(deckId, count)
	if err != nil {
		s.logger.Errorf("Error in draw %d cards from deck %s", count, remaining)
		return nil, err
	}

	return cards, nil
}
