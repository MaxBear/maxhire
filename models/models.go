package models

import (
	"fmt"
	"time"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
	applicationspb "github.com/MaxBear/maxhire/proto/gen/go/applications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application struct {
	Date     time.Time  `json:"date"`
	Company  string     `json:"company"`
	Position string     `json:"position"`
	Status   gcp.Status `json:"status"`
}

func NewApplication(a *applicationspb.Application) *Application {
	return &Application{
		Date:    a.GetDate().AsTime(),
		Company: a.GetCompany(),
	}
}

func (application *Application) Validate() error {
	if application.Date.IsZero() {
		return fmt.Errorf("invalid date")
	}
	if application.Company == "" {
		return fmt.Errorf("invalid company name")
	}
	return nil
}

func (application *Application) Pb() *applicationspb.Application {
	res := &applicationspb.Application{
		Date:     timestamppb.New(application.Date),
		Company:  application.Company,
		Position: application.Position,
		Status:   applicationspb.StatusType(application.Status),
	}
	return res
}

func ToApplication(email *gcp.Email) *Application {
	return &Application{
		Date:     email.EmailRecord.SentTime,
		Company:  email.Company,
		Position: email.Position,
		Status:   email.Status,
	}
}
