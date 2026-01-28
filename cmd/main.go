package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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

func genApplicationData(start_time, end_time, jsonFile, csvFile string) error {
	ctx := context.Background()

	err := godotenv.Load("../configs/.env")
	if err != nil {
		log.Printf("Error loading .env file, error: %s", err.Error())
		return err
	}

	llm, err := analyzer.New()
	if err != nil {
		log.Printf("Error initialize Llm analyzers, error: %s", err.Error())
		return err
	}

	// Initialize gcp app script service
	appScriptDeploymentId := os.Getenv("APP_SCRIPT_DEPLOYMENT_ID")

	s, err := gcpAppScriptService.New(
		ctx,
		gcpAppScriptService.WithOauthRedirectPort(8080),
		gcpAppScriptService.WithOauthRedirectUrl("http://localhost:8080"),
		gcpAppScriptService.WithCredFile("../configs/gcp_app_script_credentials.json"),
		gcpAppScriptService.WithTokFile("../configs/gcp_oauth_token.json"),
		gcpAppScriptService.WithAppScriptDeploymentId(appScriptDeploymentId),
		gcpAppScriptService.WithLlmAnalyzer(llm),
	)
	if err != nil {
		log.Printf("error initializing GCP App Script Service, err :%s", err.Error())
		return err
	}

	emails, err := s.GetApplicationEmails(start_time, end_time)
	if err != nil {
		log.Printf("Error get application emails using Gcp App Script service, error: %s", err.Error())
		return err
	}
	// Also print as JSON for easy copying
	jsonOutput, _ := json.MarshalIndent(emails, "", "  ")
	fmt.Printf("\nJSON output:\n%s\n", string(jsonOutput))

	applications, errs := s.ParseApplicationEmails(emails)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("Error analyze application emails, error: %s", err.Error())
		}
		return err
	}

	log.Printf("Success! Parsed %d email records:\n", len(emails))

	if len(applications) > 0 {
		applications.Print()
		applications.ToCsv(csvFile)
		applications.ToJson(jsonFile)
	}

	return nil
}

func main() {
	csv := flag.String("csv", "applications.csv", "csv file contains job application data")
	json := flag.String("json", "applications.json", "csv file contains job application data")
	load := flag.Bool("load", false, "load existing job application data from csv file")
	gen := flag.Bool("gen", false, "load existing job application data from csv file")
	start_time := flag.String("start_time", "", "start time for filtering email application confirmations, format: 2006-01-01")
	end_time := flag.String("end_time", "", "end time for filtering email application confirmations, format: 2006-01-02")
	flag.Parse()

	if *load {
		apps, err := gcp.FromJson(*json)
		if err != nil {
			os.Exit(1)
		}
		apps.Print()
		os.Exit(0)
	}

	if *gen {
		if !validTimeRange(*start_time, *end_time) {
			log.Printf("Invalid time range")
			os.Exit(1)
		}

		err := genApplicationData(*start_time, *end_time, *json, *csv)
		if err != nil {
			os.Exit(1)
		}

		os.Exit(0)
	}
}
