package sheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type credentials struct {
	PrivateKey   string `json:"private_key"`
	PrivateKeyID string `json:"private_key_id"`
	Email        string `json:"client_email"`
	TokenURL     string `json:"token_uri"`
}

var service *sheets.Service

func init() {
	Authorize()
}

// Authorize bot service account in google sheets
func Authorize() {
	credentials, err := readConfig(os.Getenv("G_CRED_FILE"))
	if err != nil {
		log.Fatalf("Unable to read credentials")
	}

	service, err = getSheetsClient(credentials)
	if err != nil {
		log.Fatalf("Unable to connect to google sheets")
	}
}

// Check if module works
func Check() {
	// Change the Spreadsheet Id with yours
	spreadsheetID := os.Getenv("SHEET_ID")

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

func readConfig(fileName string) (*credentials, error) {
	creds, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	var config credentials
	err = json.Unmarshal(creds, &config)
	if err != nil {
		log.Fatalf("Unable to unmarshall credentials: %v", err)
	}

	return &config, err
}

func getSheetsClient(cred *credentials) (*sheets.Service, error) {
	// Create a JWT configurations object for the Google service account
	config := &jwt.Config{
		Email:        (*cred).Email,
		PrivateKey:   []byte((*cred).PrivateKey),
		PrivateKeyID: (*cred).PrivateKeyID,
		TokenURL:     (*cred).TokenURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	client := config.Client(oauth2.NoContext)

	// Create a service object for Google sheets
	service, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return service, err
}
