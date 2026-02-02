package server

import (
	"context"

	"github.com/MaxBear/maxhire/models"
	applicationspb "github.com/MaxBear/maxhire/proto/gen/go/applications/v1"
	"github.com/MaxBear/maxhire/service"
)

type Server struct {
	service service.Service

	applicationspb.UnimplementedApplicationsServer
}

func New(svc service.Service) *Server {
	return &Server{
		service: svc,
	}
}

func (i *Server) ListApplications(ctx context.Context, req *applicationspb.ListApplicationsRequest) (*applicationspb.ApplicationsResponse, error) {
	applications, err := i.service.ListApplications(ctx)
	if err != nil {
		return nil, err
	}

	pbApplications := make([]*applicationspb.Application, len(applications))

	for i, application := range applications {
		pbApplications[i] = application.Pb()
	}

	return &applicationspb.ApplicationsResponse{
		Applications: pbApplications,
	}, nil
}

func (i *Server) SetApplications(ctx context.Context, req *applicationspb.SetApplicationsRequest) (*applicationspb.ApplicationsResponse, error) {
	applications := make([]*models.Application, 0)

	for _, application := range req.Applications {
		a := models.NewApplication(application)
		applications = append(applications, a)
	}

	if err := i.service.SetApplications(ctx, applications); err != nil {
		return nil, err
	}

	return &applicationspb.ApplicationsResponse{
		Applications: req.Applications,
	}, nil
}
