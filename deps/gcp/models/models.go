package models

import (
	"encoding/json"
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

type EmailRecords []EmailRecord

func (ut *EmailRecords) UnmarshalJSON(dat []byte) error {
	var data []EmailRecord
	err := json.Unmarshal(dat, &data)
	if err != nil {
		return err
	}
	*ut = data
	return nil
}
