package database

import "fmt"

// Pause contents question, user and his answer
type Pause struct {
	ID       int
	User     int
	Question int
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
		INSERT INTO pauses("user", question)
		VALUES($1, $2)
		RETURNING id;
	`, p.User, p.Question).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("Unable to save pause in DB: %v", err)
	}

	return nil
}
