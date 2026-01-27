package models

import (
	"testing"

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
		"Thanks for applying, Nancy!",
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
