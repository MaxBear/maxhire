package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) error {
	return godotenv.Load("../../../configs/.env")
}

func TestCompany(t *testing.T) {
	err := setup(t)
	require.Nil(t, err)

	tcs := []string{
		"Jane",
		"Jane Doe",
		"Senior Data Engineer, Core Experience",
		"Senior Go Backend Engineer",
		"Senior Software Development Engineer - ML Platform",
		"Senior Software Engineer (Full Stack, Backend-leaning)",
		"Senior Software Engineer (golang) - Poker",
		"Senior Software Engineer (Kubernetes), Systems",
		"Senior Software Engineer II - Observability",
		"Senior Software Engineer II- Developer Tooling Experience",
		"Senior Software Engineer Python (Django)- Vancouver",
		"Senior Software Engineer, Core Experience",
		"Senior Software Engineer, Data Pipelines",
		"Senior Software Engineer, Front-End",
		"Software Engineer – Developer Workflows & Infrastructure Automation",
		"Sr. Software Engineer",
		"Thank you for applying",
		"Thank You For Applying!",
		"Thank You for Your Application",
		"Thank you for your application",
		"Thank you for your application!",
		"Thank you for your interest!",
		"Thanks for applying, Jane!",
		"Thank you for your application!",
	}

	for _, tc := range tcs {
		c := Company(tc)
		assert.Equal(t, true, c.Invalid())
	}
}

func TestSender(t *testing.T) {
	tcs := []string{
		"no-reply@dropbox.com",
	}

	for _, tc := range tcs {
		s := Sender(tc)
		company, noreply := s.Domain()
		assert.Equal(t, true, noreply)
		assert.Equal(t, "dropbox.com", company)
	}
}

func TestFromJson(t *testing.T) {
	tcs := []struct {
		apps     Emails
		jsonFile string
	}{
		{
			apps: []*Email{
				{
					Company: "Samsara",
					Status:  Pending,
					EmailRecord: &RawEmailRecord{
						SentTime:   time.Date(2025, time.October, 3, 19, 02, 06, 0, time.UTC),
						Subject:    "Thank you for applying to Samsara",
						FullSender: "no-reply@us.greenhouse-mail.io",
						Domain:     "us.greenhouse-mail.io",
					},
				},
			},
			jsonFile: "../../../test_data/applications1.json",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			apps, err := FromJson(tc.jsonFile)
			require.Nil(t, err)
			assert.Equal(t, tc.apps, apps)
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	tcs := []struct {
		apps     Emails
		expected Emails
	}{
		{
			apps: []*Email{
				{
					Company:  "Pinterest",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 21, 15, 34, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Pending,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.January, 28, 2, 39, 5, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 20, 16, 15, 1, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Pending,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 10, 14, 31, 5, 0, time.UTC),
					},
				},
			},
			expected: []*Email{
				{
					Company:  "Pinterest",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 21, 15, 34, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Applied,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.January, 28, 2, 39, 5, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 20, 16, 15, 1, 0, time.UTC),
					},
				},
				{
					Company:  "Pinterest",
					Status:   Applied,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 10, 14, 31, 5, 0, time.UTC),
					},
				},
			},
		},
		{
			apps: []*Email{
				{
					Company:  "Twilio",
					Status:   Pending,
					Position: "Senior Software Engineer, Platform Observability",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 15, 50, 45, 0, time.UTC),
					},
				},
				{
					Company:  "Twilio",
					Status:   Reject,
					Position: "Senior Software Engineer L3",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 25, 15, 59, 52, 0, time.UTC),
					},
				},
			},
			expected: []*Email{
				{
					Company:  "Twilio",
					Status:   Pending,
					Position: "Senior Software Engineer, Platform Observability",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 15, 50, 45, 0, time.UTC),
					},
				},
				{
					Company:  "Twilio",
					Status:   Reject,
					Position: "Senior Software Engineer L3",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 25, 15, 59, 52, 0, time.UTC),
					},
				},
			},
		},
		{
			apps: []*Email{
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Full Stack Engineer, Money as a Service",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 15, 32, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.December, 1, 18, 13, 1, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Backend Engineer, Payments and Risk",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 28, 20, 45, 5, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 26, 14, 3, 38, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Technical Operations, Integration Reliability Engineer, Link",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 17, 15, 33, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Backend Engineer/API, Payments and Risk",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.October, 29, 14, 54, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Backend Engineer, Data",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.October, 3, 19, 24, 6, 0, time.UTC),
					},
				},
			},
			expected: []*Email{
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Full Stack Engineer, Money as a Service",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2026, time.February, 4, 15, 32, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.December, 1, 18, 13, 1, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Applied,
					Position: "Backend Engineer, Payments and Risk",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 28, 20, 45, 5, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Reject,
					Position: "not specified",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 26, 14, 3, 38, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Applied,
					Position: "Technical Operations, Integration Reliability Engineer, Link",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.November, 17, 15, 33, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Backend Engineer/API, Payments and Risk",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.October, 29, 14, 54, 10, 0, time.UTC),
					},
				},
				{
					Company:  "Stripe",
					Status:   Pending,
					Position: "Backend Engineer, Data",
					EmailRecord: &RawEmailRecord{
						SentTime: time.Date(2025, time.October, 3, 19, 24, 6, 0, time.UTC),
					},
				},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			tc.apps.UpdateStatus()
			for _, app := range tc.apps {
				fmt.Printf("%+v\n", app)
			}
			assert.Equal(t, tc.expected, tc.apps)
		})
	}
}
