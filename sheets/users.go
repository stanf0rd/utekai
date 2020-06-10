package sheets

import (
	"fmt"
	"log"
	"os"

	"github.com/stanf0rd/utekai/database"
	"google.golang.org/api/sheets/v4"
)

var usersPageName = os.Getenv("USERS_PAGE_NAME")

// Check if module works
func Check() {
	// Define the Sheet Name and fields to select
	readRange := "Лист1!A2:B"

	// Pull the data from the sheet
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Display pulled data
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Printf("%s, %s\n", row[0], row[1])
		}
	}
}

// PrintUsers pushes all users to gsheet
func PrintUsers(users []database.User) error {
	var vr sheets.ValueRange
	for _, user := range users {
		vr.Values = append(vr.Values, []interface{}{
			user.ID, user.TelegramID, user.Anonymous,
		})
	}

	if err := fillPage(vr, spreadsheetID, usersPageName); err != nil {
		return fmt.Errorf("Cannot fill users page: %v", err)
	}

	return nil
}

// AddUserToSheet prints user in the first empty row
func AddUserToSheet(user database.User) error {
	filledCount, err := getFilledRowsCount(spreadsheetID, usersPageName)
	if err != nil {
		return fmt.Errorf("Unable to count row count %v", err)
	}

	firstEmpty := filledCount + 1
	writeRange := fmt.Sprintf("%s!A%d:C", usersPageName, firstEmpty)

	vr := sheets.ValueRange{
		Values: [][]interface{}{
			{user.ID, user.TelegramID, user.Anonymous},
		},
	}

	_, err = service.Spreadsheets.Values.Update(
		spreadsheetID, writeRange, &vr,
	).ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("Unable to add user to sheet: %v", err)
	}

	return nil
}
