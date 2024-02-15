package gitAPI

import (
	"fmt"
	"net/http"
	"strings"
)

// CheckRun refers to GitHub's checkRun instance.
// This should only be created using NewCheckRun method.
type CheckRun struct {
	URL   string
	Token string
}

type CheckRunOutput struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// NewCheckRun creates a GitHub checkRun on a specified target.
// Note: Created checkRun's status will be in_progress.
func NewCheckRun(url string, head string, token string) *CheckRun {

	type createCheckRunBody struct {
		Name    string         `json:"name"`
		HeadSHA string         `json:"head_sha"`
		Status  string         `json:"status"`
		Output  CheckRunOutput `json:"output"`
	}

	type createCheckRunResponse struct {
		ID int `json:"id"`
	}

	requestURL := strings.ReplaceAll(url, "/hooks", "/check-runs")

	var response createCheckRunResponse

	err := RESTRequest(
		GitRequestOptions{
			Method: http.MethodPost, URL: requestURL, Body: createCheckRunBody{
				Name:    "Bobby",
				HeadSHA: head,
				Status:  "in_progress",
				Output: CheckRunOutput{
					Title: "Building in progress",
					Summary: "Build server is building your project" +
						"\nFor more information visit https://bobby.pryter.me/task_id/log",
				},
			}, Token: token,
		}, &response,
	)

	if err != nil {
		return nil
	}

	resUrl := fmt.Sprintf("%s/%d", requestURL, response.ID)

	checkrun := &CheckRun{URL: resUrl, Token: token}

	return checkrun
}

type CheckRunConclusion string

var (
	ConclusionSuccess   CheckRunConclusion = "success"
	ConclusionFailure   CheckRunConclusion = "failure"
	ConclusionCancelled CheckRunConclusion = "cancelled"
	ConclusionTimedOut  CheckRunConclusion = "timed_out"
)

func (c *CheckRun) Update(
	status string,
	conclusion CheckRunConclusion,
	output CheckRunOutput,
) {
	type updateCheckRunBody struct {
		Status     string             `json:"status"`
		Conclusion CheckRunConclusion `json:"conclusion"`
		Output     CheckRunOutput     `json:"output"`
	}

	type updateCheckRunResponse struct {
		id string
	}

	var response updateCheckRunResponse
	_ = RESTRequest(
		GitRequestOptions{
			Method: http.MethodPatch, URL: c.URL, Body: updateCheckRunBody{
				Status:     status,
				Conclusion: conclusion,
				Output:     output,
			}, Token: c.Token,
		}, response,
	)

}
