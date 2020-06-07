package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres DB driver
)

var db *sql.DB

func init() {
	host := os.Getenv("PG_APP_HOST")
	user := os.Getenv("PG_APP_USER")
	password := os.Getenv("PG_APP_PASSWORD")
	dbname := os.Getenv("PG_APP_DB")

	connStr := fmt.Sprintf(`
		host=%v
		user=%v
		password=%v
		dbname=%v
		sslmode=disable
	`, host, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// User - basic user type
type User struct {
	ID         int
	TelegramID int
	Anonymous  bool
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

// Save creates and saves new user in database
// writes userID in struct
func (u *User) Save() error {
	err := db.QueryRow(`
		INSERT INTO "users"("telegramID", anonymous)
		VALUES($1, $2)
		RETURNING id;
	`, u.TelegramID, u.Anonymous).Scan(&u.ID)

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

// GetAllUsers returns all users collected in database
// func GetAllUsers() []User {
// 	rows, err := db.Query("SELECT * FROM users")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	users := make([]User, 0)

// 	for rows.Next() {
// 		var u User
// 		if err := rows.Scan(&u); err != nil {
// 			// Check for a scan error.
// 			// Query rows will be closed with defer.
// 			log.Fatal(err)
// 		}
// 		users = append(users, u)
// 	}

// 	return users
// }
