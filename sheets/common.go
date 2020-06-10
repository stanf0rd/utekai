package sheets

import (
	"errors"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

// getFilledRowsCount reads table and returns count of filled rows
func getFilledRowsCount(spreadsheetID string, page string) (int, error) {
	readRange := fmt.Sprintf("%s!A:A", page)

	column, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	return len(column.Values), nil
}

// clearTableExceptHeader clears table, ignoring first row with column names
// returns count of cleared rows
func clearTableExceptHeader(spreadsheetID string, page string) (int, error) {
	clearRange := fmt.Sprintf("%s!A2:Z", page)

	column, err := service.Spreadsheets.Values.Get(spreadsheetID, clearRange).Do()
	if err != nil {
		return 0, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	_, err = service.Spreadsheets.Values.Clear(
		spreadsheetID, clearRange, &sheets.ClearValuesRequest{},
	).Do()

	if err != nil {
		return 0, fmt.Errorf("Unable to clear rows data: %v", err)
	}

	return len(column.Values), nil
}

func fillPage(vr sheets.ValueRange, spreadsheetID string, page string) error {
	_, err := clearTableExceptHeader(spreadsheetID, page)
	if err != nil {
		return fmt.Errorf("Cannot clear page %s: %v", page, err)
	}

	count := len(vr.Values)
	if count == 0 {
		return nil
	}

	rowWidth := len(vr.Values[0])
	if rowWidth > 25 {
		return errors.New("Too many fields in structs")
	}
	cellLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	writeRange := fmt.Sprintf(
		"%s!A2:%c%d", page, cellLetters[rowWidth-1], 1+count,
	)

	_, err = service.Spreadsheets.Values.Update(
		spreadsheetID, writeRange, &vr,
	).ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("Unable to fill page: %v", err)
	}

	return nil
}
