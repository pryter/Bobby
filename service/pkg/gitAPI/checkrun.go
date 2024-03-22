package gitAPI

import (
	"Bobby/pkg/comm"
	"fmt"
	"net/http"
	"strings"
)

// CheckRun refers to GitHub's checkRun instance.
// This should only be created using NewCheckRun method.

// NewCheckRun creates a GitHub checkRun on a specified target.
// Note: Created checkRun's status will be in_progress.
func NewCheckRun(url string, body comm.CheckRunBody, token string) (string, error) {

	type createCheckRunResponse struct {
		ID int `json:"id"`
	}

	requestURL := strings.ReplaceAll(url, "/hooks", "/check-runs")

	var response createCheckRunResponse

	err := RESTRequest(
		GitRequestOptions{
			Method: http.MethodPost, URL: requestURL, Body: body, Token: token,
		}, &response,
	)

	if err != nil {
		return "", err
	}

	resUrl := fmt.Sprintf("%s/%d", requestURL, response.ID)

	return resUrl, nil
}

func UpdateCheckRun(
	payload comm.UpdateCheckRunPayload,
	token string,
) {
	type updateCheckRunBody struct {
		Status     string                  `json:"status"`
		Conclusion comm.CheckRunConclusion `json:"conclusion"`
		Output     comm.CheckRunOutput     `json:"output"`
	}

	type updateCheckRunResponse struct {
		id string
	}

	var response updateCheckRunResponse
	_ = RESTRequest(
		GitRequestOptions{
			Method: http.MethodPatch, URL: payload.Url, Body: updateCheckRunBody{
				Status:     payload.Status,
				Conclusion: payload.Conclusion,
				Output:     payload.Output,
			}, Token: token,
		}, response,
	)

}
