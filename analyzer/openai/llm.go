package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"golang.org/x/sync/semaphore"

	gcpModels "github.com/MaxBear/maxhire/deps/gcp/models"
)

const (
	MAX_QUERIES = 5
)

type Ai struct {
	llm *openai.LLM
}

func New() (*Ai, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, err
	}
	return &Ai{
		llm: llm,
	}, nil
}

func (ai *Ai) extractApplicationDetails(ctx context.Context, message string) (status, jobTitle, companyName string, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	tool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "extract_application_details",
			Description: "Extracts the application status, job title, and company name from a job application email message",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status": map[string]any{
						"type":        "string",
						"description": "The status of the job application, either 'accept' or 'reject'",
						"enum":        []string{"accept", "reject", "pending"},
					},
					"job_title": map[string]any{
						"type":        "string",
						"description": "The job title or position name mentioned in the email",
					},
					"company_name": map[string]any{
						"type":        "string",
						"description": "The name of the company being applied to or mentioned in the email",
					},
				},
				"required": []string{"status", "job_title", "company_name"},
			},
		},
	}

	// Call the model using GenerateContent (the modern method)
	resp, err := ai.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Analyze the email message and extract: 1) whether it is an acceptance ('accept') or rejection ('reject') for a job application, 2) the job title or position name mentioned in the email, and 3) the company name."),
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}, llms.WithTools([]llms.Tool{tool}))
	if err != nil {
		log.Printf("llm error when trying to guess message details, err: %s", err.Error())
		return "", "", "", err
	}

	// Parse the extracted result from the tool calls
	if len(resp.Choices) > 0 && len(resp.Choices[0].ToolCalls) > 0 {
		var result struct {
			Status      string `json:"status"`
			JobTitle    string `json:"job_title"`
			CompanyName string `json:"company_name"`
		}
		args := resp.Choices[0].ToolCalls[0].FunctionCall.Arguments
		if err := json.Unmarshal([]byte(args), &result); err == nil {
			return result.Status, result.JobTitle, result.CompanyName, nil
		} else {
			log.Printf("llm error when trying to parse message details, err: %s", err.Error())
			return "", "", "", err
		}
	}

	return "", "", "", fmt.Errorf("llm unable to determine message details from email")
}

func (ai *Ai) AnalyzeEmails(ctx context.Context, emails gcpModels.Emails) []error {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	sem := semaphore.NewWeighted(MAX_QUERIES)
	errs := []error{}

	for i := range emails {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v\n", err)
			errs = append(errs, err)
			break
		}

		wg.Add(1)
		go func(idx int) {
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			log.Printf("analyzing email %d\n", idx)

			status, title, companyName, err := ai.extractApplicationDetails(ctx, emails[idx].EmailRecord.Msg)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			if sstatus, err := gcpModels.ParseStatus(status); err == nil {
				emails[idx].Status = sstatus
			}

			emails[idx].Position = title

			// Use company name from status guess if available and not already set
			if companyName != "" && emails[idx].Company == "" {
				emails[idx].Company = companyName
			}
		}(i)
	}
	wg.Wait()

	return errs
}
