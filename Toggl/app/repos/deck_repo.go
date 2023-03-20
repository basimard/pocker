package repos

import (
	"database/sql"
	"log"
	"strings"
	"toggl/app/config"
	"toggl/app/dtos"
	"toggl/app/models"
	"toggl/app/utils"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type DeckRepository interface {
	CreateDeck(deck models.Deck) (string, error)
	OpenDeck(deckId string) (*dtos.RespOpenDeck, error)
	CheckDeckExist(deckId string) (bool, error)
	DrawCard(deckId string, count int) (*dtos.RespDrawDeck, error)
}

type Repository struct {
	logger   *logrus.Logger
	testMode bool
	config   *config.Config
}

func setupDb(isTest bool, conf *config.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if isTest {
		db, err = sql.Open("sqlite3", conf.Database.TestPath)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	db, err = sql.Open("sqlite3", conf.Database.ProdPath)
	if err != nil {
		return nil, err
	}

	return db, nil

}

// Setup new database repository
func NewRepository(logger *logrus.Logger, testMode bool, config *config.Config) *Repository {

	// open the database
	db, err := setupDb(testMode, config)

	if err != nil {
		logger.Error(err)
	}
	defer db.Close()

	sqlStmt := `create table if not exists decks (
		id text not null primary key,
		shuffled boolean,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	  );
	  
	  create table if not exists cards (
		id text not null primary key,
		value text,
		suit text,
		deck_id text not null,
		drawn int not null DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		foreign key(deck_id) references decks(id) on delete cascade
	  );
	  
	  delete from decks;
	  delete from cards;
	  `

	// execute the SQL statements
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{logger: logger, testMode: testMode, config: config}
}

// Create deck
func (r *Repository) CreateDeck(deck *models.Deck) (string, error) {

	var deckId = utils.Generate_uuid()
	// open the database
	db, err := setupDb(r.testMode, r.config)

	if err != nil {
		r.logger.Error(err)
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		r.logger.Errorf("Error %s in begin database transaction", err)
		return "", err
	}
	defer func() {
		if err != nil {
			r.logger.Errorf("Error %s in roleback transaction", err)
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// insert new deck
	deckStmt := `
        INSERT INTO decks(id, shuffled) VALUES(?, ?);
    `

	_, err = tx.Exec(deckStmt, deckId, deck.Shuffled)
	if err != nil {
		r.logger.Errorf("Error %s in executing %s", err, deckStmt)
		return "", err
	}

	// insert cards for deck
	cardStmt := `
        INSERT INTO cards(id, value, suit, deck_id) VALUES
    `
	args := make([]interface{}, 0, 52)
	placeholders := make([]string, 0, len(deck.Cards))
	for i := 0; i < len(deck.Cards); i++ {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		args = append(args, utils.Generate_uuid(), deck.Cards[i].Value, deck.Cards[i].Suit, deckId)
	}
	cardStmt += strings.Join(placeholders, ", ")
	_, err = tx.Exec(cardStmt, args...)
	if err != nil {
		r.logger.Errorf("Error %s in executing %s", err, cardStmt)
		return "", err
	}

	return deckId, nil
}

// Open deck
func (r *Repository) OpenDeck(deckId string) (*dtos.RespOpenDeck, error) {

	db, err := setupDb(r.testMode, r.config)

	if err != nil {
		r.logger.Error(err)
	}

	defer db.Close()

	var deck dtos.RespOpenDeck
	deckQuery := `
        SELECT id, shuffled
        FROM decks
        WHERE id = ?
    `
	err = db.QueryRow(deckQuery, deckId).Scan(&deck.DeckID, &deck.Shuffled)
	if err != nil {
		r.logger.Errorf("Error %s in querying %s with %s", err, deckQuery, deckId)
		return nil, err
	}

	cardsQuery := `
        SELECT value, suit
        FROM cards
        WHERE deck_id = ? AND drawn = 0
        ORDER BY created_at
    `
	rows, err := db.Query(cardsQuery, deckId)
	deck.Remaining = 0
	if err != nil {
		r.logger.Errorf("Error %s in querying %s with %s", err, cardsQuery, deckId)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card dtos.RespOpenDeckCard
		err := rows.Scan(&card.Value, &card.Suit)
		if err != nil {
			r.logger.Errorf("Error %s in scanning %s and %s", err, "card.Value", "card.Suit")
			return nil, err
		}
		deck.Remaining += 1
		card.Code = string(card.Value[0]) + string(card.Suit[0])
		deck.Cards = append(deck.Cards, card)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error %s in scanning ", err)
		return nil, err
	}

	return &deck, nil
}

// Check is id exist
func (r *Repository) CheckDeckExist(deckId string) (bool, error) {

	db, err := setupDb(r.testMode, r.config)

	if err != nil {
		r.logger.Error(err)
	}

	defer db.Close()

	deckQuery := `
        SELECT EXISTS(
            SELECT 1 FROM decks WHERE id = ?
        )
    `
	var exist bool
	err = db.QueryRow(deckQuery, deckId).Scan(&exist)
	if err != nil {
		r.logger.Error("Error %s in querying %s with param %s", err, deckQuery, deckId)
		return false, err
	}

	return exist, nil
}

// draw cards from deck
func (r *Repository) DrawCard(deckId string, count int) (*dtos.RespDrawDeck, error) {

	db, err := setupDb(r.testMode, r.config)

	if err != nil {
		r.logger.Error(err)
	}

	defer db.Close()

	// draw cards
	cardsQuery := `
        SELECT id, value, suit
        FROM cards
        WHERE deck_id = ? AND drawn = 0
        ORDER BY created_at
        LIMIT ?
    `
	rows, err := db.Query(cardsQuery, deckId, count)
	if err != nil {
		r.logger.Errorf("Error %s in querying %s with parmas %s and %d", err, cardsQuery, deckId, count)
		return nil, err
	}
	defer rows.Close()

	var cardIds []string
	var cards []dtos.RespDrawCard
	for rows.Next() {
		var card models.Card
		err := rows.Scan(&card.Id, &card.Value, &card.Suit)
		if err != nil {
			r.logger.Errorf("Error %s in scan %s with parmas %s and %s", err, "card.Id", "card.Value", "card.Suit")
			return nil, err
		}

		cardIds = append(cardIds, card.Id)
		cards = append(cards, dtos.RespDrawCard{
			Value: card.Value,
			Suit:  card.Suit,
			Code:  string(card.Value[0]) + string(card.Suit[0]),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error %s in retriving", err)
		return nil, err
	}

	// update drawn status for cards
	updateQuery := `
        UPDATE cards SET drawn = 1 WHERE id IN (?` + strings.Repeat(",?", len(cardIds)-1) + `);
    `
	args := make([]interface{}, len(cardIds))
	for i, id := range cardIds {
		args[i] = id
	}
	_, err = db.Exec(updateQuery, args...)
	if err != nil {
		r.logger.Errorf("Error %s in updating %s with params %s", err, updateQuery, args)
		return nil, err
	}

	deck := &dtos.RespDrawDeck{
		Cards: cards,
	}
	return deck, nil
}
