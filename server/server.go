package server

import (
	"context"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
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
	filters := &service.ListApplicationsFilters{}

	// Get company filter
	company := req.GetCompany()
	if company != "" {
		filters.Company = company
	}

	// Convert status filter - only filter if explicitly set to REJECT or SUCCESS
	// (PENDING is 0, which is also the zero value, so we can't distinguish "unset" from "set to PENDING")
	status := req.GetStatus()
	if status == applicationspb.StatusType_REJECT || status == applicationspb.StatusType_SUCCESS {
		gcpStatus := gcp.Status(status)
		filters.Status = &gcpStatus
	}

	// Convert date filters
	if req.GetStartDate() != nil {
		startDate := req.GetStartDate().AsTime()
		filters.StartDate = &startDate
	}

	if req.GetEndDate() != nil {
		endDate := req.GetEndDate().AsTime()
		filters.EndDate = &endDate
	}

	applications, err := i.service.ListApplications(ctx, filters)
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
