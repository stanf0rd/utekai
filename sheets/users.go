package sheets

import (
	"fmt"
	"os"

	"github.com/stanf0rd/utekai/database"
	"google.golang.org/api/sheets/v4"
)

var usersPageName = os.Getenv("USERS_PAGE_NAME")

// PrintUsers pushes all users to gsheet
func PrintUsers(users []database.User) error {
	var vr sheets.ValueRange
	for _, user := range users {
		vr.Values = append(vr.Values, []interface{}{
			user.ID, user.TelegramID, user.Anonymous, user.Admin,
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
	writeRange := fmt.Sprintf("%s!A%d:D", usersPageName, firstEmpty)

	vr := sheets.ValueRange{
		Values: [][]interface{}{
			{user.ID, user.TelegramID, user.Anonymous, user.Admin},
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
