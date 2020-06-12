package sheets

import (
	"fmt"
	"os"

	"github.com/stanf0rd/utekai/database"
	"google.golang.org/api/sheets/v4"
)

var pausesPageName = os.Getenv("PAUSES_PAGE_NAME")

// PrintPauses pushes all pauses to gsheet
func PrintPauses(pauses []database.Pause) error {
	var vr sheets.ValueRange
	for _, p := range pauses {
		vr.Values = append(vr.Values, []interface{}{
			p.ID, p.User, p.Status, p.Question,
		})
	}

	if err := fillPage(vr, spreadsheetID, pausesPageName); err != nil {
		return fmt.Errorf("Cannot fill pauses page: %v", err)
	}

	return nil
}
