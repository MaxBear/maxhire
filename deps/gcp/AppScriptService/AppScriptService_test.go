package AppScriptService

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	analyzer "github.com/MaxBear/maxhire/analyzer/openai"
	"github.com/MaxBear/maxhire/deps/gcp/models"
)

func setupAppScriptService(t *testing.T) (*AppScriptService, error) {
	err := godotenv.Load("../../../configs/.env")
	if err != nil {
		log.Printf("Error loading .env file, error: %s", err.Error())
		return nil, err
	}

	ctx := context.Background()

	llm, err := analyzer.New()
	if err != nil {
		log.Printf("Error initialize Llm analyzers, error: %s", err.Error())
		return nil, err
	}

	return New(ctx,
		WithCredFile("../../../configs/gcp_app_script_credentials.json"),
		WithTokFile("../../../configs/gcp_oauth_token.json"),
		WithLlmAnalyzer(llm),
	)
}

func TestExtractCompany(t *testing.T) {
	service, err := setupAppScriptService(t)
	require.Nil(t, err)

	var emails models.EmailRecords
	emails = []*models.EmailRecord{
		{
			Subject:    "Update on Your Application for Sr. Software Engineer",
			FullSender: "dynata <dynata@myworkday.com>",
		},
		{
			Subject:    "Thank you for your application!",
			FullSender: "Careers lululemon <noreply-careers@lululemon.com>",
		},
	}

	expected := []string{
		"dynata",
		"Careers lululemon",
	}

	apps, errs := service.ParseApplicationEmails(emails)
	require.Empty(t, errs)
	require.Equal(t, len(expected), len(apps))
	for i, app := range apps {
		assert.Equal(t, expected[i], string(app.Company))
	}
}
