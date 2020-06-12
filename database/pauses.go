package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// Pause contents question, user and his answer
type Pause struct {
	ID           int
	User         int
	Question     int
	Status       string
	AskMessageID int
	ChatID       int64
}

// GetAskedQuestionsIDs finds all questions which were already asked to user
func GetAskedQuestionsIDs(userID int) (res []int, err error) {
	rows, err := db.Query("SELECT question FROM pauses WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("Cannot get from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ID int
		if err := rows.Scan(&ID); err != nil {
			return nil, fmt.Errorf("Cannot read question IDs from database: %v", err)
		}
		res = append(res, ID)
	}

	return res, nil
}

// Save creates and saves new pause in database
// writes pause ID in struct
func (p *Pause) Save() error {
	err := db.QueryRow(`
		INSERT INTO pauses("user", question, message_id, chat_id)
		VALUES($1, $2, $3, $4)
		RETURNING id;
	`, p.User, p.Question, p.AskMessageID, p.ChatID).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("Unable to save pause in DB: %v", err)
	}

	return nil
}

// UpdateStatus sets new status to DB
func (p *Pause) UpdateStatus() error {
	err := db.QueryRow(`
		UPDATE pauses
		SET status = $1
		WHERE id = $2
		RETURNING id;
	`, p.Status, p.ID).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("Unable to update pause status in DB: %v", err)
	}

	return nil
}

// MessageSig returns messageID and chatID
// for Editable interface from telegramBot
func (p Pause) MessageSig() (string, int64) {
	return strconv.Itoa(p.AskMessageID), int64(p.ChatID)
}

// Recipient returns user telegramID
// for Recipient interface from telegramBot
func (p Pause) Recipient() string {
	row := db.QueryRow(`
		SELECT "telegramID" FROM "users"
		WHERE id = (
			SELECT "user" FROM pauses
			WHERE id = $1
		)
		LIMIT 1;
`, p.ID)

	var userTelegramID int
	if err := row.Scan(&userTelegramID); err != nil {
		log.Fatalf("Cannot get active pause: %v", err)
	}

	return strconv.Itoa(userTelegramID)
}

// GetQuestion returns questions found by pause id
func (p Pause) GetQuestion() (*Question, error) {
	return GetQuestionByID(p.Question)
}

// GetActivePauseByUserID searches for not ended pauses without answers
func GetActivePauseByUserID(ID int) (*Pause, error) {
	row := db.QueryRow(`
		SELECT id, "user", question, status, message_id, chat_id
		FROM pauses
		WHERE "user" = $1
		AND status <> 'done'
		AND status <> 'failed'
		LIMIT 1;
	`, ID)

	var p Pause
	if err := row.Scan(
		&p.ID,
		&p.User,
		&p.Question,
		&p.Status,
		&p.AskMessageID,
		&p.ChatID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Cannot get active pause: %v", err)
		}
		return nil, nil
	}

	return &p, nil
}

// GetAllPauses returns all pauses collected in database
func GetAllPauses() ([]Pause, error) {
	rows, err := db.Query(`
		SELECT id, "user", question, status, message_id, chat_id
		FROM pauses
	`)
	if err != nil {
		return nil, fmt.Errorf("Cannot get pauses from database: %v", err)
	}
	defer rows.Close()
	pauses := make([]Pause, 0)

	for rows.Next() {
		var p Pause
		if err := rows.Scan(
			&p.ID,
			&p.User,
			&p.Question,
			&p.Status,
			&p.AskMessageID,
			&p.ChatID,
		); err != nil {
			return nil, fmt.Errorf("Cannot scan pauses from database: %v", err)
		}
		pauses = append(pauses, p)
	}

	return pauses, nil
}
