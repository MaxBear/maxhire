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
					Subject:    "Update on Your Application for Sr. Software Engineer",
					FullSender: "dynata <dynata@myworkday.com>",
					Msg:        " Dear xxx,   Thank you for your interest in the Sr. Software Engineer at Dynata. We  appreciate the time and effort you put into your application.  After careful consideration, we regret to inform you that we have decided  to move forward with other candidates whose qualifications more closely  align with the requirements of the role. While we won’t be progressing with  your application at this time, we encourage you to explore future  opportunities with us as they arise.  We sincerely appreciate your interest in Dynata and wish you the best in  your job search and future career endeavors.  Best regards, Dynata  dynata.com \u003chttps://www.dynata.com/\u003e This email was intended for nayang@maxbearwiz.com ",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Thank you for your application!",
					FullSender: "Careers lululemon <noreply-careers@lululemon.com>",
					Msg:        "",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Thank you for your interest in Uber",
					FullSender: "Do Not Reply at Uber \u003cuber+email+3m58k-5b0a3c20d7@talent.icims.com\u003e",
					Msg:        "Hi xxx,   Thank you so much for taking the time to apply for the Sr. Software Engineer role at Uber. We’ve reviewed your application and, while your resume and skills are impressive, the hiring team has decided to move forward with other candidates whose skills and experience are a closer match to what the team needs right now.   That said, we think you have great skills, and we’d love to stay connected! We encourage you to visit our Uber Careers website to explore other opportunities that may be a better fit, now or in the future.   If you have any questions please don’t hesitate to reach out and we wish you all the best in your job search.   Sincerely, Do Not Reply at Uber    Uber takes your privacy seriously. Please note by agreeing to engage with Uber Recruiting, you are consenting on the collection and use of your candidate information. To learn how we collect and use your information, and learn about your rights under GDPR, please visit our Privacy Statement.                           Get help  Privacy  Terms     Community                       Uber Technologies  1725 3rd Street, San Francisco 94158 Uber.com                        ************************************************** This message was sent to nayang@maxbearwiz.com. If you don't want to receive these emails from this company in the future, please go to: https://tracking.icims.com/f/a/7VPi6FtxobLzms8Dqz4oDQ~~/AAIB5hA~/fdm4QYAZ2HZ6D6szeIAf5WGZrQ97M8q4ccn2-HnZ25H-tw8sX8iJsrcXi6DRdnIclo4fV7V301N8VLgqRdv1wgcu-G2eh-qtHXpASNjKcgFokE8lGLnrQ1FyvcfekyRALOkBZBkhsTJiVe8e_m4KCFDcbVO2nu4FVLqqIYaeaqI~  /n \u0026#169; Uber, 1455 Market Street San Francisco CA 94103 USA/n ",
				},
			},
			{
				EmailRecord: &gcpModels.RawEmailRecord{
					Subject:    "Thank you for your interest in Instacart!",
					FullSender: "no-reply@instacart.com",
					Msg:        "Hi xxx,  Thanks for your interest in Instacart! We have received your application for our Senior Site Reliability Engineer II  role and are working as fast as we can to review it. If our hiring team sees a match for what they’re looking for, we’ll be in touch about next steps, stay tuned.  In the meantime, get a taste of what life is like at Instacart!  * See what our team is up to on theBlog ( https://news.instacart.com/ ) andLinkedIn ( https://www.linkedin.com/company/instacart )  * Meet some of our “Carrots” and check out our SF office onThe Muse ( https://www.themuse.com/profiles/instacart )  * Learn about what ourengineers are building ( https://tech.instacart.com/ ) or what ourdesign team is pear-fecting ( https://medium.com/instacart-design )  * Get a taste of our culture on our Taste of Instacart Blog ( https://www.instacart.com/company/taste-of-instacart/ )  * Check out what our team says about Instacart being agreat place to work ( https://www.greatplacetowork.com/certified-company/7014078 )  Here's what to expect from our interview process:  We'll be in touch,  The Instacart Recruiting Team  ----  At Instacart, we use BrightHire ( https://brighthire.com/ ) to record our interviews as part of our efforts to ensure an engaging hiring experience (more conversation, less note-taking). Your privacy is very important to us. If you prefer not to have any future interview recorded, you can opt-out by clicking here ( https://app.brighthire.ai/candidate-opt-out/f6sN6Ua0TqefpX7fN1XD-g ). This decision will not affect your candidacy with us in any way. Please note that we cannot provide access to interview recordings or transcripts.  Brighthire is just one of many ways we've adopted generative AI into how we work smarter. Learn more about other ways we are using AI internally ( https://tech.instacart.com/unlocking-efficiency-how-ava-became-our-ai-productivity-partner-f1a560686361 ) and in our Instacart App ( https://www.instacart.com/company/updates/bringing-inspirational-ai-powered-search-to-the-instacart-app-with-ask-instacart/ ).",
				},
			},
		},
		res: []*gcpModels.Email{
			{
				Company:  "dynata",
				Status:   gcpModels.Reject,
				Position: "Sr. Software Engineer",
			},
			{
				Company:  "Careers lululemon",
				Status:   gcpModels.Pending,
				Position: "",
			},
			{
				Company:  "Uber",
				Status:   gcpModels.Reject,
				Position: "Sr. Software Engineer",
			},
			{
				Company:  "Instacart",
				Status:   gcpModels.Pending,
				Position: "Senior Site Reliability Engineer II",
			},
		},
	}

	ai, ctx := setup(t)

	err := ai.AnalyzeEmails(ctx, tc.in)
	assert.Nil(t, err)
	for i, out := range tc.in {
		assert.Equal(t, tc.res[i].Company, out.Company)
		assert.Equal(t, tc.res[i].Position, out.Position)
		assert.Equal(t, tc.res[i].Status, out.Status)
	}
}
