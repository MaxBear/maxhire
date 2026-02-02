package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
	"github.com/MaxBear/maxhire/models"
)

type Service interface {
	SetApplications(context.Context, []*models.Application) error
	ListApplications(context.Context) ([]*models.Application, error)
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

func (s *serviceImpl) ListApplications(ctx context.Context) ([]*models.Application, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.applications, nil
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
