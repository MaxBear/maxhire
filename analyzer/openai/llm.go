package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
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

func (ai *Ai) GuessCompanyAppliedBasedOnSubject(ctx context.Context, subject string) (string, error) {
	return ai.guessCompany(ctx, "subject", subject)
}

func (ai *Ai) GuessCompanyAppliedBasedOnSender(ctx context.Context, sender string) (string, error) {
	return ai.guessCompany(ctx, "sender", sender)
}
