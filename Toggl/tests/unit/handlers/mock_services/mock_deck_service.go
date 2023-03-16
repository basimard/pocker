package mock_services

import (
	"toggl/app/dtos"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

// MockDeckService is a mock implementation of the DeckService interface
type MockDeckService struct {
	logger *logrus.Logger
	ctrl   *gomock.Controller
}

// NewMockDeckService creates a new mock of the DeckService interface
func NewMockDeckService(logger *logrus.Logger, ctrl *gomock.Controller) *MockDeckService {
	return &MockDeckService{
		logger: logger,
		ctrl:   ctrl,
	}
}

// CreateNewDeck is a mock implementation of the CreateNewDeck method
func (m *MockDeckService) CreateNewDeck(shuffle bool, cards string) (*dtos.RespCreateDeck, error) {
	ret := m.ctrl.Call(m, "CreateNewDeck", shuffle, cards)
	return ret[0].(*dtos.RespCreateDeck), nil
}

// EXPECTCreateNewDeck is a helper method for configuring expectations for the CreateNewDeck method
func (m *MockDeckService) ExpectCreateNewDeck(shuffle bool, cards string, deck *dtos.RespCreateDeck, err error) *gomock.Call {
	return m.ctrl.RecordCall(m, "CreateNewDeck", shuffle, cards).Return(deck, err)
}

// OpenDeck is a mock implementation of the OpenDeck method
func (m *MockDeckService) OpenDeck(deckId string) (*dtos.RespOpenDeck, error) {
	ret := m.ctrl.Call(m, "OpenDeck", deckId)
	return ret[0].(*dtos.RespOpenDeck), nil
}

// ExpectOpenDeck is a helper method for configuring expectations for the OpenDeck method
func (m *MockDeckService) ExpectOpenDeck(deckId string, resp *dtos.RespOpenDeck, err error) *gomock.Call {
	return m.ctrl.RecordCall(m, "OpenDeck", deckId).Return(resp, err)
}

// DrawCard is a mock implementation of the DrawCard method
func (m *MockDeckService) DrawCard(deckId string, count int) (*dtos.RespDrawDeck, error) {
	ret := m.ctrl.Call(m, "DrawCard", deckId, count)
	return ret[0].(*dtos.RespDrawDeck), nil
}

// EXPECTDrawCard is a helper method for configuring expectations for the DrawCard method
func (m *MockDeckService) ExpectDrawCard(deckId string, count int, resp *dtos.RespDrawDeck, err error) *gomock.Call {
	return m.ctrl.RecordCall(m, "DrawCard", deckId, count).Return(resp, err)
}
