package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	gcpAppScriptService "github.com/MaxBear/maxhire/deps/gcp/AppScriptService"
)

func main() {
	ctx := context.Background()

	s, err := gcpAppScriptService.New(
		ctx,
		gcpAppScriptService.WithCredFile("credentials.json"),
		gcpAppScriptService.WithOauthRedirectPort(8080),
		gcpAppScriptService.WithOauthRedirectUrl("http://localhost:8080"),
		gcpAppScriptService.WithTokFile("token.json"),
		gcpAppScriptService.WithAppScriptDeploymentId("---"),
	)

	if err != nil {
		log.Panic(err.Error())
	}

	emails, err := s.GetApplicationEmails(
		"2026-01-01",
		"2026-01-07",
	)

	if err != nil {
		log.Panic(err.Error())
	}

	// Successfully parsed
	fmt.Printf("Success! Parsed %d email records:\n", len(emails))
	for i, record := range emails {
		fmt.Printf("\n[%d] Subject: %s\n", i+1, record.Subject)
		fmt.Printf("%10s: %s\n", "SentTime", record.SentTime.Local().Format(time.RFC1123Z))
		fmt.Printf("%10s: %s\n", "Sender", record.FullSender)
		fmt.Printf("%10s: %s\n", "Domain", record.Domain)
	}
	// Also print as JSON for easy copying
	jsonOutput, _ := json.MarshalIndent(emails, "", "  ")
	fmt.Printf("\nJSON output:\n%s\n", string(jsonOutput))
}
