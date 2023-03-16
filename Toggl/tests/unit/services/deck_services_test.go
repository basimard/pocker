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
