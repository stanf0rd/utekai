package sheets

import (
	"log"
	"os"
)

var spreadsheetID = os.Getenv("SHEET_ID")

func init() {
	authorize()
	log.Println("Authorized to Google services")
}
