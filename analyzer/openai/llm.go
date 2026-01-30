package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	gcpModels "github.com/MaxBear/maxhire/deps/gcp/models"
)

var queryMsgType = map[string]bool{
	"subject": true,
	"sender":  true,
}

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

func (ai *Ai) guessCompany(ctx context.Context, message_type, message string) (string, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	var res string

	_, ok := queryMsgType[message_type]
	if !ok {
		return res, fmt.Errorf("invalid message type %q", message_type)
	}

	tool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "extract_company",
			Description: fmt.Sprintf("Extracts the company name from email %s", message_type),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"company_name": map[string]any{
						"type":        "string",
						"description": "The name of the company being applied to",
					},
				},
				"required": []string{"company_name"},
			},
		},
	}

	// Cal the model using GenerateContent (the modern method)
	resp, err := ai.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, fmt.Sprintf("Extract only the company name from email %s provided.", message_type)),
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}, llms.WithTools([]llms.Tool{tool}))
	if err != nil {
		log.Printf("llm error when try to guess company applied based on %s, err : %s", message_type, err.Error())
		return res, err
	}

	// 4. Parse the extracted result from the tool calls
	if len(resp.Choices) > 0 && len(resp.Choices[0].ToolCalls) > 0 {
		var result struct {
			CompanyName string `json:"company_name"`
		}
		args := resp.Choices[0].ToolCalls[0].FunctionCall.Arguments
		if err := json.Unmarshal([]byte(args), &result); err == nil {
			res = result.CompanyName
			return res, nil
		} else {
			log.Printf("llm error when try to guess company applied based on %s, err : %s", message_type, err.Error())
			return res, err
		}
	}

	return res, fmt.Errorf("llm unable to guess company name based on sender")
}

func (ai *Ai) guessApplicationStatus(ctx context.Context, message string) (status, jobTitle string, err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer cancelFunc()

	tool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "extract_application_details",
			Description: "Extracts the application status and job title from a job application email message",
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
				},
				"required": []string{"status", "job_title"},
			},
		},
	}

	// Call the model using GenerateContent (the modern method)
	resp, err := ai.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Analyze the email message and extract: 1) whether it is an acceptance ('accept') or rejection ('reject') for a job application, and 2) the job title or position name mentioned in the email."),
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}, llms.WithTools([]llms.Tool{tool}))
	if err != nil {
		log.Printf("llm error when trying to guess message details, err: %s", err.Error())
		return "", "", err
	}

	// Parse the extracted result from the tool calls
	if len(resp.Choices) > 0 && len(resp.Choices[0].ToolCalls) > 0 {
		var result struct {
			Status   string `json:"status"`
			JobTitle string `json:"job_title"`
		}
		args := resp.Choices[0].ToolCalls[0].FunctionCall.Arguments
		if err := json.Unmarshal([]byte(args), &result); err == nil {
			return result.Status, result.JobTitle, nil
		} else {
			log.Printf("llm error when trying to parse message details, err: %s", err.Error())
			return "", "", err
		}
	}

	return "", "", fmt.Errorf("llm unable to determine message details from email")
}

func (ai *Ai) guessCompanyName(ctx context.Context, email *gcpModels.RawEmailRecord) (string, error) {
	companyName, err := ai.guessCompany(ctx, "subject", email.Subject)
	if err != nil {
		log.Printf("llm error when try to guess company name based on email subject %q, err : %s", email.Subject, err.Error())
		return "", err
	}
	if !gcpModels.Company(companyName).Invalid() {
		return companyName, nil
	}
	ddomain := gcpModels.Sender(email.FullSender)
	companyName, noreply := ddomain.Domain()
	if noreply {
		return companyName, nil
	}
	companyName, err = ai.guessCompany(ctx, "sender", string(email.FullSender))
	if err != nil {
		log.Printf("llm error when try to guess company name based on sender address %q, err : %s", email.FullSender, err.Error())
		return "", err
	}
	return companyName, nil
}

func (ai *Ai) AnalyzeEmails(ctx context.Context, emails gcpModels.Emails) []error {
	errs := []error{}
	for i, email := range emails {
		log.Printf("analyzing email %d\n", i)

		companyName, err := ai.guessCompanyName(ctx, email.EmailRecord)
		if err != nil {
			errs = append(errs,
				fmt.Errorf("failed to analyze email %q from %s sent at %s",
					email.EmailRecord.Subject,
					email.EmailRecord.FullSender,
					email.EmailRecord.SentTime.Format(time.RFC1123Z)))
			continue
		}
		email.Company = companyName

		if len(email.EmailRecord.Msg) == 0 {
			email.Status = gcpModels.Pending
			continue
		}

		status, title, err := ai.guessApplicationStatus(ctx, email.EmailRecord.Msg)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if sstatus, err := gcpModels.ParseStatus(status); err == nil {
			email.Status = sstatus
		}

		email.Position = title
	}
	return errs
}
