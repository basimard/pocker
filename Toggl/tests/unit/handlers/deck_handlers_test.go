package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"toggl/app/dtos"
	"toggl/app/handlers"
	"toggl/tests/unit/handlers/mock_services"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewDeckHandlerWithoutParamsReturnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up expected inputs and outputs for the CreateNewDeck method
	expectedDeck := &dtos.RespCreateDeck{
		DeckID:    "a251071b-662f-44b6-ba11-e24863039c59",
		Shuffled:  false,
		Remaining: 50,
	}

	expectedErr := errors.New("some error")
	mockDeckService.ExpectCreateNewDeck(false, "", expectedDeck, expectedErr)

	// Set up the HTTP request and response
	req, errs := http.NewRequest("POST", "/v1/create-deck", nil)
	if errs != nil {
		fmt.Println(errs)
	}
	w := httptest.NewRecorder()

	// Call the handler
	handler.CreateNewDeckHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":50}`
	actual := w.Body.String()
	assert.Equal(t, actual, expected)
}

func TestCreateDeckHandlerWithCardsParmsReturnSuccess(t *testing.T) {
	var cards = "AS,2S"
	var shuffled = "true"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up expected inputs and outputs for the CreateNewDeck method
	expectedDeck := &dtos.RespCreateDeck{
		DeckID:    "a251071b-662f-44b6-ba11-e24863039c59",
		Shuffled:  true,
		Remaining: 2,
	}

	expectedErr := errors.New("some error")
	mockDeckService.ExpectCreateNewDeck(true, cards, expectedDeck, expectedErr)
	req, err := http.NewRequest("POST", "/v1/create-deck?cards="+cards+"&shuffle="+shuffled, nil)
	assert.NoError(t, err)

	resRec := httptest.NewRecorder()

	handler.CreateNewDeckHandler(resRec, req)

	assert.Equal(t, http.StatusOK, resRec.Code)

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":2}`
	actual := resRec.Body.String()
	assert.Equal(t, expected, actual)

}

func TestOpenDeckHandlerWithParamsReturnSuccess(t *testing.T) {
	var id = `a251071b-662f-44b6-ba11-e24863039c59`
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	expectedDeck := &dtos.RespOpenDeck{
		DeckID:    id,
		Shuffled:  true,
		Remaining: 2,
		Cards: []dtos.RespOpenDeckCard{
			{Code: "AC", Value: "ACE", Suit: "SPADES"},
			{Code: "1H", Value: "10", Suit: "HEARTS"},
		},
	}

	expectedErr := errors.New("some error")
	mockDeckService.ExpectOpenDeck(id, expectedDeck, expectedErr)
	req, err := http.NewRequest("GET", "/open-deck?deck_id="+id, nil)
	assert.NoError(t, err)

	resRec := httptest.NewRecorder()

	handler.OpenDeckHandler(resRec, req)

	assert.Equal(t, http.StatusOK, resRec.Code)

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":2,"cards":[{"code":"AC","value":"ACE","suit":"SPADES"},{"code":"1H","value":"10","suit":"HEARTS"}]}`
	actual := resRec.Body.String()
	assert.Equal(t, expected, actual)

}

func TestOpenDeckHandlerWithEmptyParamsReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up the HTTP request and response with an empty deck_id
	req, _ := http.NewRequest("GET", "/v1/deck/open-deck?deck_id=", nil)
	w := httptest.NewRecorder()

	expectedDeck := &dtos.RespOpenDeck{}

	expectedErr := errors.New("some error")
	// Expect that the service layer is not called
	mockDeckService.ExpectOpenDeck("", expectedDeck, expectedErr).Times(0)

	// Call the handler
	handler.OpenDeckHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := strings.TrimSpace(w.Body.String())
	actual := "Deck id parameter is required"
	assert.Equal(t, expected, actual)
}

func TestOpenDeckHandlerWithInvalidParamsReturnsError(t *testing.T) {

	id := `invalid_id`
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up the HTTP request and response with an empty deck_id
	req, _ := http.NewRequest("GET", "/v1/deck/open-deck?deck_id="+id, nil)
	w := httptest.NewRecorder()

	expectedDeck := &dtos.RespOpenDeck{}

	expectedErr := errors.New("some error")
	// Expect that the service layer is not called
	mockDeckService.ExpectOpenDeck("", expectedDeck, expectedErr).Times(0)

	// Call the handler
	handler.OpenDeckHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := strings.TrimSpace(w.Body.String())
	actual := "Invalid deck id"
	assert.Equal(t, expected, actual)
}

func TestDrawCardHandlerWithParamsReturnSuccess(t *testing.T) {
	var id = `a251071b-662f-44b6-ba11-e24863039c59`
	var count = 3
	logger := logrus.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up expected inputs and outputs for the CreateNewDeck method
	expectedDeck := &dtos.RespDrawDeck{
		Cards: []dtos.RespDrawCard{
			{Code: "AC", Value: "ACE", Suit: "SPADES"},
			{Code: "1H", Value: "10", Suit: "HEARTS"},
		},
	}

	expectedErr := errors.New("some error")
	mockDeckService.ExpectDrawCard(id, count, expectedDeck, expectedErr)

	req, err := http.NewRequest("GET", "/v1/draw-cards?deck_id="+id+"&count="+fmt.Sprintf("%d", count), nil)
	assert.NoError(t, err)

	resRec := httptest.NewRecorder()

	handler.DrawCardHandler(resRec, req)

	assert.Equal(t, http.StatusOK, resRec.Code)

	expected := `{"cards":[{"code":"AC","value":"ACE","suit":"SPADES"},{"code":"1H","value":"10","suit":"HEARTS"}]}`
	actual := resRec.Body.String()
	assert.Equal(t, expected, actual)
}

func TestDrawCardHandlerWithEmptyDeckIdParamsReturnError(t *testing.T) {
	var id = ``
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up the HTTP request and response with an empty deck_id
	req, _ := http.NewRequest("GET", "/v1/deck/draw-cards?deck_id="+id, nil)
	w := httptest.NewRecorder()

	expectedDeck := &dtos.RespDrawDeck{}

	expectedErr := errors.New("some error")
	// Expect that the service layer is not called
	mockDeckService.ExpectDrawCard("", 3, expectedDeck, expectedErr).Times(0)

	// Call the handler
	handler.DrawCardHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := strings.TrimSpace(w.Body.String())
	actual := "Deck id parameter is required"
	assert.Equal(t, expected, actual)

}

func TestDrawCardHandlerWithInvalidDeckIdParamReturnError(t *testing.T) {
	var id = `invalid`
	var count = `3`
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up the HTTP request and response with an empty deck_id
	req, _ := http.NewRequest("GET", "/v1/deck/draw-cards?deck_id="+id+"count="+count, nil)
	w := httptest.NewRecorder()

	expectedDeck := &dtos.RespDrawDeck{}

	expectedErr := errors.New("some error")
	// Expect that the service layer is not called
	mockDeckService.ExpectDrawCard(id, 3, expectedDeck, expectedErr).Times(0)

	// Call the handler
	handler.DrawCardHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := strings.TrimSpace(w.Body.String())
	actual := "Invalid deck id"
	assert.Equal(t, expected, actual)

}

func TestDrawCardHandlerWithInvalidCountParamReturnError(t *testing.T) {
	var id = "a251071b-662f-44b6-ba11-e24863039c59"
	var count = "A"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logrus.New()
	mockDeckService := mock_services.NewMockDeckService(logger, ctrl)

	handler := handlers.NewDeckHandler(mockDeckService, logger)

	// Set up the HTTP request and response with an empty deck_id
	req, _ := http.NewRequest("GET", "/v1/deck/draw_card?deck_id="+id+"&count="+count, nil)
	w := httptest.NewRecorder()

	expectedDeck := &dtos.RespDrawDeck{}

	expectedErr := errors.New("some error")
	// Expect that the service layer is not called
	mockDeckService.ExpectDrawCard(id, 0, expectedDeck, expectedErr).Times(0)

	// Call the handler
	handler.DrawCardHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := strings.TrimSpace(w.Body.String())
	actual := "Count parameter must be a positive integer"
	assert.Equal(t, expected, actual)

}
