package sheets

import (
	"fmt"
	"os"

	"github.com/stanf0rd/utekai/database"
	"google.golang.org/api/sheets/v4"
)

var questionsPageName = os.Getenv("QUESTIONS_PAGE_NAME")

// PrintQuestions pushes all questions to gsheet
func PrintQuestions(questions []database.Question) error {
	var vr sheets.ValueRange
	for _, q := range questions {
		vr.Values = append(vr.Values, []interface{}{
			q.ID, q.Order, q.Body,
		})
	}

	if err := fillPage(vr, spreadsheetID, questionsPageName); err != nil {
		return fmt.Errorf("Cannot fill questions page: %v", err)
	}

	return nil
}
