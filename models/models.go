package models

import (
	"fmt"
	"time"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
	applicationspb "github.com/MaxBear/maxhire/proto/gen/go/applications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application struct {
	Date       time.Time   `json:"date"`
	Company    string      `json:"company"`
	Position   string      `json:"position"`
	Status     gcp.Status  `json:"status"`
	Interviews []Interview `json:"interviews"`
}

type InterviewType int

const (
	Unspecified      InterviewType = iota // 0
	RecruiterScreen                       // 1
	ManagerScreen                         // 2
	TechCoding                            // 3
	TechSystemDesign                      // 4
	TeamMatch                             // 5
)

// String method for general printing (fmt.Println)
func (t InterviewType) String() string {
	return [...]string{
		"Unspecified",
		"RecruiterScreen",
		"ManagerScreen",
		"TechCoding",
		"TechSystemDesign",
		"TeamMatch"}[t]
}

func ParseInterviewType(s string) (InterviewType, error) {
	statusMap := map[string]InterviewType{
		"Unspecified":      Unspecified,
		"RecruiterScreen":  RecruiterScreen,
		"ManagerScreen":    ManagerScreen,
		"TechCoding":       TechCoding,
		"TechSystemDesign": TechSystemDesign,
		"TeamMatch":        TeamMatch,
	}

	if val, ok := statusMap[s]; ok {
		return val, nil
	}
	return Unspecified, fmt.Errorf("invalid status: %s", s)
}

type Interview struct {
	DateTime      time.Time     `json:"dateTime"`
	InterviewType InterviewType `json:"type"`
	DurationMin   int32         `json:"durationMin"`
}

func NewApplication(a *applicationspb.Application) *Application {
	interviews := make([]Interview, 0, len(a.GetInterviews()))
	for _, pbInterview := range a.GetInterviews() {
		interviews = append(interviews, InterviewFromPb(pbInterview))
	}
	return &Application{
		Date:       a.GetDate().AsTime(),
		Company:    a.GetCompany(),
		Position:   a.GetPosition(),
		Status:     gcp.Status(a.GetStatus()),
		Interviews: interviews,
	}
}

func InterviewFromPb(pb *applicationspb.Interview) Interview {
	durationMin := pb.GetDurationMin()
	// Default to 15 minutes if not specified (zero value)
	if durationMin == 0 {
		durationMin = 15
	}
	return Interview{
		DateTime:      pb.GetDatetime().AsTime(),
		InterviewType: InterviewType(pb.GetInterviewType()),
		DurationMin:   durationMin,
	}
}

func (i *Interview) Pb() *applicationspb.Interview {
	return &applicationspb.Interview{
		Datetime:      timestamppb.New(i.DateTime),
		InterviewType: applicationspb.InterviewType(i.InterviewType),
		DurationMin:   i.DurationMin,
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
	interviews := make([]*applicationspb.Interview, 0, len(application.Interviews))
	for _, interview := range application.Interviews {
		interviews = append(interviews, interview.Pb())
	}
	res := &applicationspb.Application{
		Date:       timestamppb.New(application.Date),
		Company:    application.Company,
		Position:   application.Position,
		Status:     applicationspb.StatusType(application.Status),
		Interviews: interviews,
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
