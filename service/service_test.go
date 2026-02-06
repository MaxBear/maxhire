package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gcp "github.com/MaxBear/maxhire/deps/gcp/models"
	"github.com/MaxBear/maxhire/models"
)

func TestSetInterviews_Success(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Create a test application
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testCompany := "TestCompany"
	application := &models.Application{
		Date:       testDate,
		Company:    testCompany,
		Position:   "Software Engineer",
		Status:     gcp.Pending,
		Interviews: []models.Interview{},
	}

	// Add the application to the service
	err = svc.SetApplications(ctx, []*models.Application{application})
	require.NoError(t, err)

	// Create interviews to set
	interview1 := &models.Interview{
		DateTime:      time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
		InterviewType: models.RecruiterScreen,
		DurationMin:   30,
	}
	interview2 := &models.Interview{
		DateTime:      time.Date(2024, 1, 25, 15, 0, 0, 0, time.UTC),
		InterviewType: models.TechCoding,
		DurationMin:   60,
	}
	interviews := []*models.Interview{interview1, interview2}

	// Set interviews
	result, err := svc.SetInterviews(ctx, testDate, testCompany, interviews)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the interviews were set
	assert.Equal(t, testDate, result.Date)
	assert.Equal(t, testCompany, result.Company)
	assert.Len(t, result.Interviews, 2)
	assert.Equal(t, interview1.DateTime, result.Interviews[0].DateTime)
	assert.Equal(t, interview1.InterviewType, result.Interviews[0].InterviewType)
	assert.Equal(t, interview1.DurationMin, result.Interviews[0].DurationMin)
	assert.Equal(t, interview2.DateTime, result.Interviews[1].DateTime)
	assert.Equal(t, interview2.InterviewType, result.Interviews[1].InterviewType)
	assert.Equal(t, interview2.DurationMin, result.Interviews[1].DurationMin)
}

func TestSetInterviews_ApplicationNotFound(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Try to set interviews for a non-existent application
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testCompany := "NonExistentCompany"
	interviews := []*models.Interview{
		{
			DateTime:      time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
			InterviewType: models.RecruiterScreen,
			DurationMin:   30,
		},
	}

	result, err := svc.SetInterviews(ctx, testDate, testCompany, interviews)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "application not found")
}

func TestSetInterviews_ReplaceExistingInterviews(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Create a test application with existing interviews
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testCompany := "TestCompany"
	application := &models.Application{
		Date:     testDate,
		Company:  testCompany,
		Position: "Software Engineer",
		Status:   gcp.Pending,
		Interviews: []models.Interview{
			{
				DateTime:      time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
				InterviewType: models.ManagerScreen,
				DurationMin:   45,
			},
		},
	}

	// Add the application to the service
	err = svc.SetApplications(ctx, []*models.Application{application})
	require.NoError(t, err)

	// Create new interviews to replace the existing ones
	newInterviews := []*models.Interview{
		{
			DateTime:      time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
			InterviewType: models.RecruiterScreen,
			DurationMin:   30,
		},
		{
			DateTime:      time.Date(2024, 1, 25, 15, 0, 0, 0, time.UTC),
			InterviewType: models.TechCoding,
			DurationMin:   60,
		},
		{
			DateTime:      time.Date(2024, 1, 30, 16, 0, 0, 0, time.UTC),
			InterviewType: models.TechSystemDesign,
			DurationMin:   90,
		},
	}

	// Set interviews (should replace existing)
	result, err := svc.SetInterviews(ctx, testDate, testCompany, newInterviews)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the old interview was replaced
	assert.Len(t, result.Interviews, 3)
	assert.Equal(t, models.RecruiterScreen, result.Interviews[0].InterviewType)
	assert.Equal(t, int32(30), result.Interviews[0].DurationMin)
	assert.Equal(t, models.TechCoding, result.Interviews[1].InterviewType)
	assert.Equal(t, int32(60), result.Interviews[1].DurationMin)
	assert.Equal(t, models.TechSystemDesign, result.Interviews[2].InterviewType)
	assert.Equal(t, int32(90), result.Interviews[2].DurationMin)
	// Verify the old interview is gone
	for _, interview := range result.Interviews {
		assert.NotEqual(t, models.ManagerScreen, interview.InterviewType)
		assert.NotEqual(t, int32(45), interview.DurationMin)
	}
}

func TestSetInterviews_EmptyInterviewsList(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Create a test application with existing interviews
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testCompany := "TestCompany"
	application := &models.Application{
		Date:     testDate,
		Company:  testCompany,
		Position: "Software Engineer",
		Status:   gcp.Pending,
		Interviews: []models.Interview{
			{
				DateTime:      time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
				InterviewType: models.ManagerScreen,
				DurationMin:   45,
			},
		},
	}

	// Add the application to the service
	err = svc.SetApplications(ctx, []*models.Application{application})
	require.NoError(t, err)

	// Set empty interviews list
	result, err := svc.SetInterviews(ctx, testDate, testCompany, []*models.Interview{})
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the interviews list is now empty
	assert.Len(t, result.Interviews, 0)
}

func TestSetInterviews_MultipleApplicationsSameCompany(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Create two applications with same company but different dates
	testCompany := "TestCompany"
	date1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC)

	application1 := &models.Application{
		Date:       date1,
		Company:    testCompany,
		Position:   "Software Engineer",
		Status:     gcp.Pending,
		Interviews: []models.Interview{},
	}
	application2 := &models.Application{
		Date:       date2,
		Company:    testCompany,
		Position:   "Senior Software Engineer",
		Status:     gcp.Pending,
		Interviews: []models.Interview{},
	}

	// Add both applications
	err = svc.SetApplications(ctx, []*models.Application{application1, application2})
	require.NoError(t, err)

	// Set interviews for the first application
	interviews1 := []*models.Interview{
		{
			DateTime:      time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
			InterviewType: models.RecruiterScreen,
			DurationMin:   30,
		},
	}

	result1, err := svc.SetInterviews(ctx, date1, testCompany, interviews1)
	require.NoError(t, err)
	assert.Len(t, result1.Interviews, 1)
	assert.Equal(t, "Software Engineer", result1.Position)
	assert.Equal(t, int32(30), result1.Interviews[0].DurationMin)

	// Set interviews for the second application
	interviews2 := []*models.Interview{
		{
			DateTime:      time.Date(2024, 2, 20, 14, 0, 0, 0, time.UTC),
			InterviewType: models.TechCoding,
			DurationMin:   60,
		},
		{
			DateTime:      time.Date(2024, 2, 25, 15, 0, 0, 0, time.UTC),
			InterviewType: models.TechSystemDesign,
			DurationMin:   90,
		},
	}

	result2, err := svc.SetInterviews(ctx, date2, testCompany, interviews2)
	require.NoError(t, err)
	assert.Len(t, result2.Interviews, 2)
	assert.Equal(t, "Senior Software Engineer", result2.Position)

	// Verify the first application's interviews weren't affected
	allApps, err := svc.ListApplications(ctx, nil)
	require.NoError(t, err)
	for _, app := range allApps {
		if app.Date.Equal(date1) && app.Company == testCompany {
			assert.Len(t, app.Interviews, 1)
			assert.Equal(t, models.RecruiterScreen, app.Interviews[0].InterviewType)
			assert.Equal(t, int32(30), app.Interviews[0].DurationMin)
		}
		if app.Date.Equal(date2) && app.Company == testCompany {
			assert.Len(t, app.Interviews, 2)
			assert.Equal(t, int32(60), app.Interviews[0].DurationMin)
			assert.Equal(t, int32(90), app.Interviews[1].DurationMin)
		}
	}
}

func TestSetInterviews_AllInterviewTypes(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, "")
	require.NoError(t, err)

	// Create a test application
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testCompany := "TestCompany"
	application := &models.Application{
		Date:       testDate,
		Company:    testCompany,
		Position:   "Software Engineer",
		Status:     gcp.Pending,
		Interviews: []models.Interview{},
	}

	// Add the application
	err = svc.SetApplications(ctx, []*models.Application{application})
	require.NoError(t, err)

	// Create interviews with all interview types
	interviews := []*models.Interview{
		{
			DateTime:      time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
			InterviewType: models.Unspecified,
			DurationMin:   0,
		},
		{
			DateTime:      time.Date(2024, 1, 21, 11, 0, 0, 0, time.UTC),
			InterviewType: models.RecruiterScreen,
			DurationMin:   30,
		},
		{
			DateTime:      time.Date(2024, 1, 22, 12, 0, 0, 0, time.UTC),
			InterviewType: models.ManagerScreen,
			DurationMin:   45,
		},
		{
			DateTime:      time.Date(2024, 1, 23, 13, 0, 0, 0, time.UTC),
			InterviewType: models.TechCoding,
			DurationMin:   60,
		},
		{
			DateTime:      time.Date(2024, 1, 24, 14, 0, 0, 0, time.UTC),
			InterviewType: models.TechSystemDesign,
			DurationMin:   90,
		},
		{
			DateTime:      time.Date(2024, 1, 25, 15, 0, 0, 0, time.UTC),
			InterviewType: models.TeamMatch,
			DurationMin:   30,
		},
	}

	// Set interviews
	result, err := svc.SetInterviews(ctx, testDate, testCompany, interviews)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify all interview types are present
	assert.Len(t, result.Interviews, 6)
	interviewTypes := make(map[models.InterviewType]bool)
	durationMap := make(map[models.InterviewType]int32)
	for _, interview := range result.Interviews {
		interviewTypes[interview.InterviewType] = true
		durationMap[interview.InterviewType] = interview.DurationMin
	}
	assert.True(t, interviewTypes[models.Unspecified])
	assert.Equal(t, int32(0), durationMap[models.Unspecified])
	assert.True(t, interviewTypes[models.RecruiterScreen])
	assert.Equal(t, int32(30), durationMap[models.RecruiterScreen])
	assert.True(t, interviewTypes[models.ManagerScreen])
	assert.Equal(t, int32(45), durationMap[models.ManagerScreen])
	assert.True(t, interviewTypes[models.TechCoding])
	assert.Equal(t, int32(60), durationMap[models.TechCoding])
	assert.True(t, interviewTypes[models.TechSystemDesign])
	assert.Equal(t, int32(90), durationMap[models.TechSystemDesign])
	assert.True(t, interviewTypes[models.TeamMatch])
	assert.Equal(t, int32(30), durationMap[models.TeamMatch])
}
