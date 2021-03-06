package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

// User - basic user type
type User struct {
	ID         int
	TelegramID int
	Anonymous  bool
	Admin      bool
}

// Exists checks if user row is already in database
// if exists, writes userID in struct
func (u *User) Exists() (bool, error) {
	err := db.QueryRow(`
		SELECT id FROM "users"
		WHERE "telegramID" = $1;
	`, u.TelegramID).Scan(&u.ID)

	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	} else {
		return false, fmt.Errorf("Unable to get user from DB: %v", err)
	}
}

// GetUserByTelegramID gets user from DB
func GetUserByTelegramID(ID int) (*User, error) {
	var u User
	err := db.QueryRow(`
		SELECT id, "telegramID", anonymous, admin FROM "users"
		WHERE "telegramID" = $1;
	`, ID).Scan(&u.ID, &u.TelegramID, &u.Anonymous, &u.Admin)

	if err != nil {
		return nil, fmt.Errorf("Unable to get user from DB: %v", err)
	}

	return &u, nil
}

// Save creates and saves new user in database
// writes userID in struct
func (u *User) Save() error {
	err := db.QueryRow(`
		INSERT INTO "users"("telegramID", anonymous, admin)
		VALUES($1, $2, $3)
		RETURNING id;
	`, u.TelegramID, u.Anonymous, u.Admin).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("Unable to save user in DB: %v", err)
	}

	return nil
}

// UpdateAnonymity updates user anonymity choise in database
// writes userID in struct
func (u *User) UpdateAnonymity() error {
	err := db.QueryRow(`
		UPDATE "users"
		SET anonymous = $2
		WHERE "telegramID" = $1
		RETURNING id;
	`, u.TelegramID, u.Anonymous).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("Unable to update user anonymity in DB: %v", err)
	}

	return nil
}

// Recipient returns TelegramID
func (u User) Recipient() string {
	return strconv.Itoa(u.TelegramID)
}

// GetAllUsers returns all users collected in database
func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT \"id\", \"telegramID\", anonymous, admin FROM \"users\"")
	if err != nil {
		return nil, fmt.Errorf("Cannot get users from database: %v", err)
	}
	defer rows.Close()
	users := make([]User, 0)

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.TelegramID, &u.Anonymous, &u.Admin); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			return nil, fmt.Errorf("Cannot read users from database: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}
