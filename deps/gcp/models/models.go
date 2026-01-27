package models

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type ApiResp struct {
	AtType string      `json:"@type"`
	Result interface{} `json:"result"` // Use interface{} to handle both string and array
}

// EmailRecord represents a single email record returned from the script
type EmailRecord struct {
	SentTime   time.Time `json:"SentTime"`
	Subject    string    `json:"Subject"`
	FullSender string    `json:"FullSender"`
	Domain     string    `json:"Domain"`
}

type EmailRecords []*EmailRecord

func (ut *EmailRecords) UnmarshalJSON(dat []byte) error {
	var data []*EmailRecord
	err := json.Unmarshal(dat, &data)
	if err != nil {
		return err
	}
	*ut = data
	return nil
}

func (in EmailRecords) ToApplications() []*Application {
	res := []*Application{}
	for _, email := range in {
		res = append(res, &Application{
			EmailRecord: email,
		})
	}
	return res
}

type Company string

func (c Company) Invalid() bool {
	substrings := []string{
		fmt.Sprintf("%s", os.Getenv("APPLICANT_FIRST_NAME")),
		fmt.Sprintf("%s!", os.Getenv("APPLICANT_LAST_NAME")),
		"jane",
		"jane!",
		"senior",
		"engineer",
		"thank you",
		"application",
		"applying",
		"your company",
		"sentaur",
		"interest",
		"infra",
	}

	found := false
	for _, sub := range substrings {
		cc := strings.ToLower(string(c))
		pattern := `\b` + regexp.QuoteMeta(sub) + `\b`
		re := regexp.MustCompile(pattern)

		if re.MatchString(cc) {
			found = true
			break
		}
	}

	return found
}

type Sender string

func (s Sender) Domain() (string, bool) {
	prefixes := []string{
		"no-reply@",
		"gh-no-reply@",
	}
	for _, prefix := range prefixes {
		domain, found := strings.CutPrefix(string(s), prefix)
		if found {
			return domain, true
		}
	}
	return "", false
}

type Status int

const (
	StatusPending Status = iota // 0
	StatusReject                // 1
	StatusSuccess               // 2
)

type Application struct {
	Company     string
	Status      Status
	EmailRecord *EmailRecord
}

type Applications []*Application

func (in Applications) ToCsv(csvFile string) error {
	f, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache application csv file, error : %v", err)
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write Header
	writer.Write([]string{"SentTime", "Subject", "FullSender", "Domain", "Company"})

	// Write Records
	for _, app := range in {
		record := []string{
			app.EmailRecord.SentTime.Format(time.RFC1123),
			app.EmailRecord.Subject,
			app.EmailRecord.FullSender,
			app.EmailRecord.Domain,
			app.Company,
		}
		writer.Write(record)
	}

	log.Printf("successfully saved applications to %s\n", csvFile)
	return nil
}

func FromCsv(csvFile string) (Applications, error) {
	var apps Applications
	apps = []*Application{}

	f, err := os.Open(csvFile)
	if err != nil {
		return apps, nil
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, _ := reader.Read()
	colMap := make(map[string]int)

	for i, name := range header {
		colMap[name] = i
	}

	// Accessing data by column name
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Use the map to get the correct index for "Email"
		t, err := time.Parse(time.RFC1123, record[colMap["SentTime"]])
		if err != nil {
			log.Printf("error parsing sent time, error: %s", err.Error())
			return apps, err
		}

		app := &Application{
			Company: record[colMap["Company"]],
			EmailRecord: &EmailRecord{
				SentTime:   t,
				Subject:    record[colMap["Subject"]],
				FullSender: record[colMap["FullSender"]],
				Domain:     record[colMap["Domain"]],
			},
		}

		apps = append(apps, app)
	}

	return apps, nil
}

func (in Applications) Print() {
	for i, app := range in {
		fmt.Printf("\n[%d] Subject: %s\n", i+1, app.EmailRecord.Subject)
		fmt.Printf("%10s: %s\n", "SentTime", app.EmailRecord.SentTime.Local().Format(time.RFC1123Z))
		fmt.Printf("%10s: %s\n", "Sender", app.EmailRecord.FullSender)
		fmt.Printf("%10s: %s\n", "Domain", app.EmailRecord.Domain)
		fmt.Printf("%10s: %s\n", "Company", app.Company)
	}
}
