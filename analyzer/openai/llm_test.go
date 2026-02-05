package analyzer

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gcpModels "github.com/MaxBear/maxhire/deps/gcp/models"
)

func setup(t *testing.T) (*Ai, context.Context) {
	godotenv.Load("../../configs/.env")
	ai, err := New()
	require.Nil(t, err)

	ctx := context.Background()

	return ai, ctx
}

func TestAnalyzeEmails(t *testing.T) {
	tc := struct {
		in  gcpModels.Emails
		res gcpModels.Emails
	}{
		in: []*gcpModels.Email{
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Thank you for your application for Software Engineer – Developer Workflows & Infrastructure Automation",
					FullSender: "no-reply@us.greenhouse-mail.io",
					Msg:        "Hello xx,  Thank you for your interest in Lyft! We wanted to let you know we received your application for Software Engineer – Developer Workflows & Infrastructure Automation, and we are delighted that you would consider joining our team. The Recruiting team will review your application and will be in touch if your qualifications match our needs for the role. If you are not selected for this position, keep an eye on our Careers ( https://www.lyft.com/careers ) page as we're growing and adding openings.  Due to the volume of applications, we want to be respectful of your time and let you know that we will only be reaching out to candidates who are an immediate fit for a role.  Want to learn more about what it's like to work for Lyft? Check out Lyft's LinkedIn page ( https://www.linkedin.com/company/lyft/life/ ) and the Lyft Blog ( https://blog.lyft.com/ ).  Thank you, The Lyft Team",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Zapier: Thanks for applying — here’s what’s next!",
					FullSender: "Zapier Hiring Team <no-reply@ashbyhq.com>",
					Msg:        "Hi xx,  Beep, boop, bop! This is our friendly Zapbot confirming we received your application for the Sr. Software Engineer (L4) role at Zapier.  Thanks for applying! We're receiving a very high volume of applications, so while a real human will review yours, it may take up to 14 days for us to get back to you with a decision. If you have any questions while we review it, our https://jobs-page-22fc3d.zapier.app/ (built with Zapier) is standing by with answers about the role, our hiring process, and life at Zapier.  While you wait, dive into these quick guides to get hands-on with Zapier and AI models. Familiarity with automation and AI will set you up for success during interviews and on the job!   - https://zapier.com/blog/get-started-with-zapier/ – build your first Zap in minutes and see our platform in action.   - https://zapier.com/blog/zapier-ai-orchestration-platform/?utm_source=chatgpt.com – shows how to plug powerful AI tools into thousands of apps and build smart, time-saving workflows in minutes.   - https://www.anthropic.com/ai-fluency – a beginner-friendly course on when to let AI help, how to ask it clearly, spot mistakes, and use it responsibly.  Talk soon!  The Zapier Team",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Update on Your Application with Affinity.co",
					FullSender: "no-reply@affinity.co",
					Msg:        "Hello xx,  Thank you for your interest in Affinity ( https://www.affinity.co/ ) and for taking the time to apply for the Senior Software Engineer, Data position. We are always amazed at how many wonderful people want to work with us, and we enjoyed reviewing your professional experience.  We are reaching out with an update on our end. The interview process advanced, and the Senior Software Engineer, Data position is now filled. We greatly appreciate you applying and sharing your skills and experience with us. We'll keep your information on file and look forward to reaching out to you when the right opportunity arises.  Thanks again for the time you invested in completing our application, and we wish you success with your job search and in your career.  Warmly,  Affinity Talent Team",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Thanks for applying to Stripe!",
					FullSender: "no-reply@stripe.com",
					Msg:        "Hi xx,  Thanks so much for submitting your application for the Backend Engineer, Data role! Our Recruiting team will review what you've sent as soon as possible. We do our best to get back to everyone who applies, but please note that, due to the number of applications we receive, we aren't always able to reply to every candidate.  In the meantime, feel free to check out our quick guide to culture at Stripe [0], our blog [1], and our LinkedIn page [2] to get a feel for what we're up to at Stripe!  Wishing you the best in your search,  Stripe  [0] https://stripe.com/jobs/culture  [1] https://stripe.com/blog  [2] https://www.linkedin.com/company/stripe  ** This email is sent from an outbound-only address, and we can't see replies.",
				},
			},
		},
		res: []*gcpModels.Email{
			{
				Company:  "Lyft",
				Status:   gcpModels.Pending,
				Position: "Software Engineer – Developer Workflows & Infrastructure Automation",
			},
			{
				Company:  "Zapier",
				Status:   gcpModels.Pending,
				Position: "Sr. Software Engineer (L4)",
			},
			{
				Company:  "Affinity",
				Status:   gcpModels.Reject,
				Position: "Senior Software Engineer, Data",
			},
			{
				Company:  "Stripe",
				Status:   gcpModels.Pending,
				Position: "Backend Engineer, Data",
			},
		},
	}

	ai, ctx := setup(t)

	err := ai.AnalyzeEmails(ctx, tc.in)
	assert.Equal(t, 0, len(err))
	for i, out := range tc.in {
		assert.Equal(t, tc.res[i].Company, out.Company)
		assert.Equal(t, tc.res[i].Position, out.Position)
		assert.Equal(t, tc.res[i].Status, out.Status)
	}
}
