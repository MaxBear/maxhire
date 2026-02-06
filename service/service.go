package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
	"github.com/MaxBear/maxhire/models"
)

type Service interface {
	SetApplications(context.Context, []*models.Application) error
	ListApplications(context.Context, *ListApplicationsFilters) ([]*models.Application, error)
	SetInterviews(context.Context, time.Time, string, []*models.Interview) (*models.Application, error)
}

type ListApplicationsFilters struct {
	Status    *gcp.Status
	Company   string
	StartDate *time.Time
	EndDate   *time.Time
}

func NewService(ctx context.Context, jsonFile string) (*serviceImpl, error) {
	applications := []*models.Application{}

	if len(jsonFile) > 0 {
		emails, err := gcp.FromJson(jsonFile)
		if err != nil {
			log.Printf("failed to load application records from %s, error: %s", jsonFile, err.Error())
			return nil, err
		}
		for _, email := range emails {
			applications = append(applications, models.ToApplication(email))
		}
		log.Printf("Successfully loaded %d applications from %s", len(emails), jsonFile)
	}

	return &serviceImpl{
		ctx:          ctx,
		applications: applications,
	}, nil
}

type serviceImpl struct {
	mu           sync.RWMutex
	applications []*models.Application
	ctx          context.Context
}

func (s *serviceImpl) ListApplications(ctx context.Context, filters *ListApplicationsFilters) ([]*models.Application, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if filters == nil {
		return s.applications, nil
	}

	var filtered []*models.Application
	for _, app := range s.applications {
		// Filter by company
		if filters.Company != "" && app.Company != filters.Company {
			continue
		}

		// Filter by status
		if filters.Status != nil && app.Status != *filters.Status {
			continue
		}

		// Filter by start date
		if filters.StartDate != nil && app.Date.Before(*filters.StartDate) {
			continue
		}

		// Filter by end date
		if filters.EndDate != nil && app.Date.After(*filters.EndDate) {
			continue
		}

		filtered = append(filtered, app)
	}

	return filtered, nil
}

func (s *serviceImpl) SetApplications(ctx context.Context, applications []*models.Application) error {
	for _, application := range applications {
		err := application.Validate()
		if err != nil {
			return fmt.Errorf("invalid application found %+v, error: %s", *application, err.Error())
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.applications = append(s.applications, applications...)

	return nil
}

func (s *serviceImpl) SetInterviews(ctx context.Context, date time.Time, company string, interviews []*models.Interview) (*models.Application, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the application by date and company
	var foundApp *models.Application
	for _, app := range s.applications {
		if app.Date.Equal(date) && app.Company == company {
			foundApp = app
			break
		}
	}

	if foundApp == nil {
		return nil, fmt.Errorf("application not found for date %v and company %s", date, company)
	}

	// Convert []*models.Interview to []models.Interview
	interviewSlice := make([]models.Interview, len(interviews))
	for i, interview := range interviews {
		interviewSlice[i] = *interview
	}

	// Set the interviews (replace existing)
	foundApp.Interviews = interviewSlice

	return foundApp, nil
}
