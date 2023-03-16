package services

import (
	"fmt"
	"os"
	"testing"
	"toggl/app/repos"
	"toggl/app/services"
	"toggl/app/utils"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	// Perform setup operations here, if any
	setup()

	// Run all the test cases
	exitCode := m.Run()

	// Perform teardown operations here, if any
	teardown()

	// Exit with the same code as the test cases
	os.Exit(exitCode)
}

func setup() {
	_, err := createTempDB()
	if err != nil {
		panic(err)
	}

}

func teardown() {
	err := os.Remove("../../../app/db/test.db")
	if err != nil {
		panic(err)
	}
}

func createTempDB() (string, error) {
	// Create a temporary file for the database
	file, err := os.Create("../../../app/db/test.db")
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer file.Close()
	dbPath := file.Name()

	return dbPath, nil
}

func TestCheckIfCreateNewDeckReturnValidDeckId(t *testing.T) {

	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, err := service.CreateNewDeck(false, "")

	// Ensure that no error was returned
	assert.NoError(t, err)

	// Parse the UUID and ensure that it is valid
	parsedUUID, err := utils.Parse_uuid(deck.DeckID)
	assert.NoError(t, err)
	assert.Equal(t, true, parsedUUID)
}

func TestCheckIfCreateNewDeckReturnValidRemainingForGivenCards(t *testing.T) {

	var sample = "AS,2S"
	var expectedCount = 2
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, err := service.CreateNewDeck(false, sample)

	// Ensure that no error was returned
	assert.NoError(t, err)

	// Parse the UUID and ensure that it is valid

	assert.NoError(t, err)
	assert.Equal(t, deck.Remaining, expectedCount)
}

func TestCheckIfCreateNewDeckReturnInValidForInvalidCards(t *testing.T) {

	var sample = "AS3,2SQ"

	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	_, err := service.CreateNewDeck(false, sample)

	// Ensure that no error was returned
	assert.EqualError(t, err, "Invalid card")

}

func TestCheckIfCreateNewDeckNonShafulledCardsSameOrder(t *testing.T) {
	var stringSample = "AS,2S,3S,4S"
	var sample = []string{"AS", "2S", "3S", "4S"}

	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, _ := service.CreateNewDeck(false, stringSample)

	deckOpend, _ := service.OpenDeck(deck.DeckID)

	for index, card := range deckOpend.Cards {

		assert.Equal(t, sample[index], card.Code)

	}

}

func TestCheckIfCreateNewDeckShafulledCardsDifferetOrder(t *testing.T) {
	var stringSample = "AS,2S"
	var sample = []string{"AS", "2S"}
	var shuffled = true
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, _ := service.CreateNewDeck(shuffled, stringSample)

	deckOpend, _ := service.OpenDeck(deck.DeckID)

	for index, card := range deckOpend.Cards {

		assert.NotEqual(t, sample[index], card.Code)

	}

}

func TestCheckIfCreateNewDeckIfCodeSuitPositionChangeReturnError(t *testing.T) {
	var sample = "SA,S2"

	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	_, err := service.CreateNewDeck(false, sample)

	// Ensure that no error was returned
	assert.EqualError(t, err, "Invalid value")

}

//More than 52 codes in query not tested

func TestCheckIfOpenDeckWithValidIdReturnCorrectData(t *testing.T) {
	var stringSample = "AS,2S"
	var sample = []string{"AS", "2S"}
	var shuffled = false
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, _ := service.CreateNewDeck(shuffled, stringSample)
	newDeckId := deck.DeckID
	deckOpend, _ := service.OpenDeck(newDeckId)

	for index, card := range deckOpend.Cards {

		assert.Equal(t, sample[index], card.Code)

	}

	assert.Equal(t, newDeckId, deckOpend.DeckID)
	assert.Equal(t, len(sample), deckOpend.Remaining)

}

func TestCheckIfOpenDeckWithNonExistIdReturnError(t *testing.T) {

	var sample = "a251071b-662f-44b6-ba11-e24863039c59"

	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	_, err := service.OpenDeck(sample)

	// Ensure that no error was returned
	assert.EqualError(t, err, "Id doesn't exist")

}

func TestCheckIfDrawCardWithValidIdReturnSuccess(t *testing.T) {

	var stringSample = "AS,2S,3S"
	var sample = []string{"AS", "2S", "3S"}
	var shuffled = false
	var count = 2
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, _ := service.CreateNewDeck(shuffled, stringSample)
	newDeckId := deck.DeckID
	drawnCards, _ := service.DrawCard(newDeckId, count)

	for index, card := range drawnCards.Cards {

		assert.Contains(t, sample[index], card.Code)
	}

	assert.Equal(t, len(drawnCards.Cards), count)

}

func TestCheckIfDrawnCardWithMoreThanRemainingCountReturnError(t *testing.T) {

	var stringSample = "AS,2S,3S"
	var shuffled = false
	var count = 4
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	// Call the CreateNewDeck method with false for shuffle and an empty string for deckID
	deck, _ := service.CreateNewDeck(shuffled, stringSample)
	newDeckId := deck.DeckID
	_, err := service.DrawCard(newDeckId, count)

	assert.EqualError(t, err, "Requested count exceeds remaining cards in deck")

}

func TestCheckIfDrawnCardsNonExisitingIdReturnError(t *testing.T) {
	var sample = "a251071b-662f-44b6-ba11-e24863039c59"
	var count = 4
	// Create a new logger
	logger := logrus.New()

	// Create a new repository in test mode
	repo := repos.NewRepository(logger, true)

	// Create a new deck service using the repository
	service := services.NewDeckService(logger, repo)

	_, err := service.DrawCard(sample, count)

	assert.EqualError(t, err, "Id doesn't exist")

}
