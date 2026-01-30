package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	analyzer "github.com/MaxBear/maxhire/analyzer/openai"
	gcpAppScriptService "github.com/MaxBear/maxhire/deps/gcp/AppScriptService"
	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
)

func validTimeRange(start_time, end_time string) bool {
	// time.DateOnly is a predefined constant for "2006-01-02"
	ts, err := time.Parse(time.DateOnly, start_time)
	if err != nil {
		log.Printf("error parsing start time, error: %s", err.Error())
		return false
	}

	te, err := time.Parse(time.DateOnly, end_time)
	if err != nil {
		log.Printf("error parsing end time, error: %s", err.Error())
		return false
	}

	if ts.Sub(te) > 24*time.Hour {
		log.Printf("Start time must be more than 24 hours after end time.")
		return false
	}

	return true
}

func genApplicationData(ctx context.Context, start_time, end_time, jsonFile, csvFile string, useLlm bool) error {
	// Initialize gcp app script service
	appScriptDeploymentId := os.Getenv("APP_SCRIPT_DEPLOYMENT_ID")

	s, err := gcpAppScriptService.New(
		ctx,
		gcpAppScriptService.WithOauthRedirectPort(8080),
		gcpAppScriptService.WithOauthRedirectUrl("http://localhost:8080"),
		gcpAppScriptService.WithCredFile("../configs/gcp_app_script_credentials.json"),
		gcpAppScriptService.WithTokFile("../configs/gcp_oauth_token.json"),
		gcpAppScriptService.WithAppScriptDeploymentId(appScriptDeploymentId),
	)
	if err != nil {
		log.Printf("error initializing GCP App Script Service, err :%s", err.Error())
		return err
	}

	raws, err := s.GetApplicationEmails(start_time, end_time)
	if err != nil {
		log.Printf("Error get application emails using Gcp App Script service, error: %s", err.Error())
		return err
	}

	jsonOutput, _ := json.MarshalIndent(raws, "", "  ")
	fmt.Printf("\nJSON output:\n%s\n", string(jsonOutput))

	emails := raws.ToEmails()

	if len(emails) > 0 {
		emails.Print()
		emails.ToCsv(csvFile)
		emails.ToJson(jsonFile)
	}

	return nil
}

func fname(orig string) string {
	extension := filepath.Ext(orig)

	// 2. Remove it from the filename
	nameWithoutExtension := orig[0 : len(orig)-len(extension)]

	return fmt.Sprintf("%s_llm", nameWithoutExtension)
}

func analyzeApplicationData(ctx context.Context, jsonFile string) error {
	emails, err := gcp.FromJson(jsonFile)
	if err != nil {
		log.Printf("unable to load application data from %s\n, error: %s", jsonFile, err.Error())
		return err
	}

	emails.Print()

	llm, err := analyzer.New()
	if err != nil {
		log.Printf("error initialize Llm analyzers, error: %s", err.Error())
		return err
	}

	errs := llm.AnalyzeEmails(ctx, emails)
	if len(errs) > 0 {
		log.Println("error analyzing email applications, errors:")
		for i, err := range errs {
			log.Printf("%d %s", i, err.Error())
		}
	}

	if len(emails) > 0 {
		emails.Print()
		emails.ToCsv(fmt.Sprintf("%s.csv", fname(jsonFile)))
		emails.ToJson(fmt.Sprintf("%s.json", fname(jsonFile)))
	}

	return nil
}

func main() {
	csv := flag.String("csv", "raw.csv", "csv file contains job application records")
	json := flag.String("json", "raw.json", "json file contains job application records")
	gen := flag.Bool("gen", false, "generating job application records to csv file")
	start_time := flag.String("start_time", "", "start time for filtering job applications, format: 2006-01-01")
	end_time := flag.String("end_time", "", "end time for filtering job applications, format: 2006-01-02")
	llm := flag.Bool("llm", false, "using LLM to analyze job applications")

	flag.Parse()

	if *gen == false && *llm == false {
		os.Exit(0)
	}

	ctx := context.Background()

	err := godotenv.Load("../configs/.env")
	if err != nil {
		log.Printf("Error loading .env file, error: %s", err.Error())
		os.Exit(1)
	}

	if *gen {
		if !validTimeRange(*start_time, *end_time) {
			log.Printf("Invalid time range")
			os.Exit(1)
		}

		err := genApplicationData(ctx, *start_time, *end_time, *json, *csv, *llm)
		if err != nil {
			os.Exit(1)
		}

		os.Exit(0)
	}

	// Use llm to populate fields such as company name, application status etc.
	if *llm {
		err := analyzeApplicationData(ctx, *json)

		if err != nil {
			os.Exit(1)
		}
	}
}
