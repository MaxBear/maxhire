package AppScriptService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/script/v1"

	"github.com/MaxBear/maxhire/deps/gcp/models"
)

type AppScriptService struct {
	ctx                       context.Context
	withCredFile              string
	withTokFile               string
	withOauthRedirectUrl      string
	withOauthRedirectPort     int
	withAppScriptDeploymentId string
	oAuthClient               *http.Client
	scriptService             *script.Service
}

type AppScriptServiceOpt func(*AppScriptService)

func WithCredFile(credFile string) AppScriptServiceOpt {
	return func(s *AppScriptService) {
		s.withCredFile = credFile
	}
}

func WithTokFile(tokFile string) AppScriptServiceOpt {
	return func(s *AppScriptService) {
		s.withTokFile = tokFile
	}
}

func WithOauthRedirectUrl(url string) AppScriptServiceOpt {
	return func(s *AppScriptService) {
		s.withOauthRedirectUrl = url
	}
}

func WithOauthRedirectPort(port int) AppScriptServiceOpt {
	return func(s *AppScriptService) {
		s.withOauthRedirectPort = port
	}
}

func WithAppScriptDeploymentId(id string) AppScriptServiceOpt {
	return func(s *AppScriptService) {
		s.withAppScriptDeploymentId = id
	}
}

func New(ctx context.Context, opts ...AppScriptServiceOpt) (*AppScriptService, error) {
	s := &AppScriptService{
		ctx: ctx,
	}

	for _, opt := range opts {
		opt(s)
	}

	b, err := os.ReadFile(s.withCredFile)
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/script.projects",
		"https://www.googleapis.com/auth/script.scriptapp",
		"https://mail.google.com/",
		"https://www.googleapis.com/auth/spreadsheets")

	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	config.RedirectURL = s.withOauthRedirectUrl

	client, err := s.getClient(config)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	s.oAuthClient = client

	srv, err := script.NewService(ctx, option.WithHTTPClient(s.oAuthClient))
	if err != nil {
		log.Printf("Unable to retrieve Script client: %v", err)
		return nil, err
	}
	s.scriptService = srv

	return s, nil
}

// Request a token from the web, then returns the retrieved token.
func (s *AppScriptService) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	codeCh := make(chan string)
	server := &http.Server{Addr: fmt.Sprintf(":%d", s.withOauthRedirectPort)}

	// Define the callback handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Auth successful! You can return to the terminal.")
			codeCh <- code
		}
	})

	go server.ListenAndServe()

	// Wait for the code from the browser
	code := <-codeCh
	server.Shutdown(s.ctx)

	// Exchange the code for an actual Token
	tok, err := config.Exchange(s.ctx, code)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return nil, err
	}

	return tok, nil
}

// Retrieves a token from a local file.
func (s *AppScriptService) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(s.withTokFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (s *AppScriptService) saveToken(token *oauth2.Token) error {
	log.Printf("Saving credential file to: %s\n", s.withTokFile)
	f, err := os.OpenFile(s.withTokFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache oauth token, error: %v", err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func (s *AppScriptService) getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := s.tokenFromFile()
	if err == nil {
		log.Println("Using existing token. If you get AuthRequiredError, delete token.json and re-run to re-authenticate with new scopes.")
		return config.Client(s.ctx, tok), nil
	}

	log.Println("No existing token founing OAuth flow...")
	tok, err = s.getTokenFromWeb(config)
	if err != nil {
		return nil, err
	}

	if err := s.saveToken(tok); err != nil {
		return nil, err
	}

	return config.Client(s.ctx, tok), nil
}

func (s *AppScriptService) GetApplicationEmails(start_date, end_date string) (models.RawEmailRecords, error) {
	emails := []*models.RawEmailRecord{}

	req := &script.ExecutionRequest{
		Function: "runFilterMyEmails", // The name of the function in your .gs file
		Parameters: []interface{}{
			start_date,
			end_date,
		},
	}

	resp, err := s.scriptService.Scripts.Run(s.withAppScriptDeploymentId, req).Do()
	if err != nil {
		log.Printf("Unable to execute script: %v", err)
		return emails, err
	}

	if resp.Error != nil {
		log.Printf("Script error: %s\n", resp.Error.Message)
		if resp.Error.Code != 0 {
			log.Printf("Error code: %d\n", resp.Error.Code)
		}
		if len(resp.Error.Details) > 0 {
			log.Println("Error details:")
			for i, detail := range resp.Error.Details {
				detailJSON, _ := json.MarshalIndent(detail, "  ", "  ")
				log.Printf("  [%d]: %s\n", i+1, string(detailJSON))
			}
		}
		// Print full error object for debugging
		errorJSON, _ := json.MarshalIndent(resp.Error, "", "  ")
		log.Printf("\nFull error object:\n%s\n", string(errorJSON))

		return emails, fmt.Errorf("error running script, err code : %d, err :  %s",
			resp.Error.Code,
			resp.Error.Message)
	}

	var res models.ApiResp
	var parseError error

	// Strategy 1: Try to unmarshal as wrapped response
	if err := json.Unmarshal(resp.Response, &res); err == nil && res.Result != nil {
		// Check if result is a string (JSON-encoded)
		if resultStr, ok := res.Result.(string); ok {
			// It's a JSON string, unmarshal it
			if err = json.Unmarshal([]byte(resultStr), &emails); err != nil {
				log.Printf("Error parsing response: %v", parseError)
				return emails, err
			}
		}
	}

	return emails, nil
}
