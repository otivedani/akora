package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	TOKEN_FILE      string = os.Getenv("TOKEN_FILE")
	CREDENTIAL_FILE string = os.Getenv("CREDENTIAL_FILE")
	SENDER_EMAIL    string = os.Getenv("SENDER_EMAIL")
)

var config *oauth2.Config
var token *oauth2.Token

// Returns the generated client.
func GetGoogleClient(ctx context.Context) *http.Client {
	return config.Client(ctx, token)
}

func GetSenderEmail() string {
	return SENDER_EMAIL
}

func init() {
	config = getConfig()
	token = getToken()
}

func getConfig() *oauth2.Config {
	b, err := os.ReadFile(CREDENTIAL_FILE)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	scopes := []string{
		gmail.MailGoogleComScope,
		gmail.GmailSendScope,
		gmail.GmailComposeScope,
		gmail.GmailModifyScope,
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

// Retrieve a token, saves the token
func getToken() *oauth2.Token {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(TOKEN_FILE)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(TOKEN_FILE, tok)
	}
	return tok
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
