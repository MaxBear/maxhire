package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	applicationspb "github.com/MaxBear/maxhire/proto/gen/go/applications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestInterviewFromPb_DefaultDuration(t *testing.T) {
	// Test that when duration_min is 0 (not set), it defaults to 15 minutes
	pbInterview := &applicationspb.Interview{
		Datetime:      timestamppb.New(time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC)),
		InterviewType: applicationspb.InterviewType_RECRUITER_SCREEN,
		DurationMin:   0, // Not set, should default to 15
	}

	interview := InterviewFromPb(pbInterview)
	assert.Equal(t, int32(15), interview.DurationMin, "duration_min should default to 15 when not set")
	assert.Equal(t, time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC), interview.DateTime)
	assert.Equal(t, RecruiterScreen, interview.InterviewType)
}

func TestInterviewFromPb_ExplicitDuration(t *testing.T) {
	// Test that when duration_min is explicitly set, it uses that value
	pbInterview := &applicationspb.Interview{
		Datetime:      timestamppb.New(time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC)),
		InterviewType: applicationspb.InterviewType_TECH_CODING,
		DurationMin:   60, // Explicitly set to 60 minutes
	}

	interview := InterviewFromPb(pbInterview)
	assert.Equal(t, int32(60), interview.DurationMin, "duration_min should use the explicitly set value")
	assert.Equal(t, time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC), interview.DateTime)
	assert.Equal(t, TechCoding, interview.InterviewType)
}

func TestInterview_Pb(t *testing.T) {
	// Test conversion from models.Interview to protobuf
	interview := Interview{
		DateTime:      time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
		InterviewType: RecruiterScreen,
		DurationMin:   30,
	}

	pbInterview := interview.Pb()
	require.NotNil(t, pbInterview)
	assert.Equal(t, time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC), pbInterview.GetDatetime().AsTime())
	assert.Equal(t, applicationspb.InterviewType_RECRUITER_SCREEN, pbInterview.GetInterviewType())
	assert.Equal(t, int32(30), pbInterview.GetDurationMin())
}

func TestInterview_Pb_RoundTrip(t *testing.T) {
	// Test round-trip conversion: pb -> model -> pb
	originalPb := &applicationspb.Interview{
		Datetime:      timestamppb.New(time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC)),
		InterviewType: applicationspb.InterviewType_MANAGER_SCREEN,
		DurationMin:   45,
	}

	// Convert to model
	modelInterview := InterviewFromPb(originalPb)
	assert.Equal(t, int32(45), modelInterview.DurationMin)

	// Convert back to pb
	resultPb := modelInterview.Pb()
	require.NotNil(t, resultPb)
	assert.Equal(t, originalPb.GetDatetime().AsTime(), resultPb.GetDatetime().AsTime())
	assert.Equal(t, originalPb.GetInterviewType(), resultPb.GetInterviewType())
	assert.Equal(t, originalPb.GetDurationMin(), resultPb.GetDurationMin())
}

func TestInterviewFromPb_RoundTripWithDefault(t *testing.T) {
	// Test round-trip when duration_min is 0 (should default to 15)
	originalPb := &applicationspb.Interview{
		Datetime:      timestamppb.New(time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC)),
		InterviewType: applicationspb.InterviewType_TECH_SYSTEM_DESIGN,
		DurationMin:   0, // Not set
	}

	// Convert to model (should apply default of 15)
	modelInterview := InterviewFromPb(originalPb)
	assert.Equal(t, int32(15), modelInterview.DurationMin, "should default to 15 minutes")

	// Convert back to pb (will be 15, not 0)
	resultPb := modelInterview.Pb()
	require.NotNil(t, resultPb)
	assert.Equal(t, int32(15), resultPb.GetDurationMin(), "should preserve the default value")
}

func TestNewApplication_WithInterviews(t *testing.T) {
	// Test that NewApplication correctly converts interviews with default duration
	pbApplication := &applicationspb.Application{
		Date:     timestamppb.New(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
		Company:  "TestCompany",
		Position: "Software Engineer",
		Status:   applicationspb.StatusType_PENDING,
		Interviews: []*applicationspb.Interview{
			{
				Datetime:      timestamppb.New(time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC)),
				InterviewType: applicationspb.InterviewType_RECRUITER_SCREEN,
				DurationMin:   0, // Should default to 15
			},
			{
				Datetime:      timestamppb.New(time.Date(2024, 1, 25, 15, 0, 0, 0, time.UTC)),
				InterviewType: applicationspb.InterviewType_TECH_CODING,
				DurationMin:   60, // Explicitly set
			},
		},
	}

	application := NewApplication(pbApplication)
	require.NotNil(t, application)
	assert.Len(t, application.Interviews, 2)
	assert.Equal(t, int32(15), application.Interviews[0].DurationMin, "first interview should default to 15 minutes")
	assert.Equal(t, int32(60), application.Interviews[1].DurationMin, "second interview should use explicit value")
}
