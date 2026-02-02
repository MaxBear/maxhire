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
type RawEmailRecord struct {
	SentTime   time.Time `json:"SentTime"`
	Subject    string    `json:"Subject"`
	FullSender string    `json:"FullSender"`
	Domain     string    `json:"Domain"`
	Msg        string    `json:"Msg"`
}

type RawEmailRecords []*RawEmailRecord

func (ut *RawEmailRecords) UnmarshalJSON(dat []byte) error {
	var data []*RawEmailRecord
	err := json.Unmarshal(dat, &data)
	if err != nil {
		return err
	}
	*ut = data
	return nil
}

func (in RawEmailRecords) ToEmails() Emails {
	res := []*Email{}
	for _, email := range in {
		res = append(res, &Email{
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
	Pending Status = iota // 0
	Reject                // 1
	Success               // 2
)

// String method for general printing (fmt.Println)
func (s Status) String() string {
	return [...]string{"Pending", "Reject", "Success"}[s]
}

func ParseStatus(s string) (Status, error) {
	statusMap := map[string]Status{
		"pending": Pending,
		"reject":  Reject,
		"accept":  Success,
	}

	if val, ok := statusMap[s]; ok {
		return val, nil
	}
	return Pending, fmt.Errorf("invalid status: %s", s)
}

// MarshalJSON allows the enum to be marshaled as a string in JSON
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	// 1. Unmarshal the raw JSON bytes into a string
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	// 2. Map the string to the corresponding constant
	switch name {
	case "Pending":
		*s = Pending
	case "Reject":
		*s = Reject
	case "Success":
		*s = Success
	}
	return nil
}

type Email struct {
	Company     string          `json:"Company"`
	Status      Status          `json:"Status"`
	Position    string          `json:"Position"`
	EmailRecord *RawEmailRecord `json:"Email"`
}

type Emails []*Email

func (in Emails) ToJson(jsonFile string) error {
	fileData, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		log.Printf("Unable marshal json data, error : %s", err.Error())
		return err
	}

	// 2. Write to file
	err = os.WriteFile(jsonFile, fileData, 0644)
	if err != nil {
		log.Printf("Unable to cache application csv file %q, error : %v", jsonFile, err)
		return err
	}

	log.Printf("successfully saved applications to %s\n", jsonFile)
	return nil
}

func FromJson(jsonFile string) (Emails, error) {
	var emails Emails
	emails = []*Email{}

	byteValue, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Printf("Unable marshal json data, error : %s", err.Error())
		return emails, err
	}

	// 2. Unmarshal JSON into the struct
	if err := json.Unmarshal(byteValue, &emails); err != nil {
		log.Printf("Unable unmarshal application data from json file %q, error : %s", jsonFile, err.Error())
		return emails, err
	}

	return emails, nil
}

func (in Emails) ToCsv(csvFile string) error {
	f, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache application csv file, error : %v", err)
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write Header
	writer.Write([]string{"SentTime", "Subject", "FullSender", "Domain", "Company", "Position", "Status"})

	// Write Records
	for _, email := range in {
		record := []string{
			email.EmailRecord.SentTime.Format(time.RFC3339),
			email.EmailRecord.Subject,
			email.EmailRecord.FullSender,
			email.EmailRecord.Domain,
			email.Company,
			email.Position,
			email.Status.String(),
		}
		writer.Write(record)
	}

	log.Printf("successfully saved applications to %s\n", csvFile)
	return nil
}

func FromCsv(csvFile string) (Emails, error) {
	var apps Emails
	apps = []*Email{}

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

		app := &Email{
			Company: record[colMap["Company"]],
			EmailRecord: &RawEmailRecord{
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

func (in Emails) Print() {
	for i, email := range in {
		fmt.Printf("\n[%d] Subject: %s\n", i+1, email.EmailRecord.Subject)
		fmt.Printf("%10s: %s\n", "SentTime", email.EmailRecord.SentTime.Local().Format(time.RFC1123Z))
		fmt.Printf("%10s: %s\n", "Sender", email.EmailRecord.FullSender)
		fmt.Printf("%10s: %s\n", "Domain", email.EmailRecord.Domain)
		fmt.Printf("%10s: %s\n", "Msg", email.EmailRecord.Msg)
		fmt.Printf("%10s: %s\n", "Company", email.Company)
		fmt.Printf("%10s: %s\n", "Position", email.Position)
		fmt.Printf("%10s: %s\n", "Status", email.Status)
	}
}
