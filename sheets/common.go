package sheets

import (
	"fmt"
	"log"
)

func getFilledRowsCount(spreadsheetID string, page string) int {
	readRange := fmt.Sprintf("%s!A:A", page)

	column, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	log.Println(len(column.Values))

	return len(column.Values)
}
