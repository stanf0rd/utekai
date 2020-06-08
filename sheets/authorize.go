package sheets

import (
	"encoding/json"
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

// Authorize bot service account in google sheets
func authorize() {
	credentials, err := readConfig(os.Getenv("G_CRED_FILE"))
	if err != nil {
		log.Fatalf("Unable to read credentials")
	}

	service, err = getSheetsClient(credentials)
	if err != nil {
		log.Fatalf("Unable to connect to google sheets")
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
