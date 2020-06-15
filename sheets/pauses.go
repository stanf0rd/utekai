package sheets

import (
	"fmt"
	"log"
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
			p.ID, p.User, p.Status, p.Question, p.Answer,
		})
	}

	if err := fillPage(vr, spreadsheetID, pausesPageName); err != nil {
		return fmt.Errorf("Cannot fill pauses page: %v", err)
	}

	printPhotos(pauses)
	return nil
}

func printPhotos(pauses []database.Pause) {
	var vr sheets.ValueRange
	for _, p := range pauses {
		if p.Photo != "" {
			link := fmt.Sprintf("https://utekai.behind.blue%s", p.Photo)
			vr.Values = append(vr.Values, []interface{}{
				fmt.Sprintf("=HYPERLINK(\"%s\"; IMAGE(\"%s\"; 1))", link, link),
			})
		} else {
			vr.Values = append(vr.Values, []interface{}{
				"No photo",
			})
		}
	}

	count := len(vr.Values)
	if count == 0 {
		return
	}

	writeRange := fmt.Sprintf(
		"%s!F2:G%d", pausesPageName, 1+count,
	)

	_, err := service.Spreadsheets.Values.Update(
		spreadsheetID, writeRange, &vr,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		log.Printf("Unable to print photos: %v", err)
	}
}
