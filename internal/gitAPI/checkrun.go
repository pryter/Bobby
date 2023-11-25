package gitAPI

import (
	"fmt"
	"net/http"
	"strings"
)

type CheckRun struct {
	URL   string
	Token string
}

func NewCheckRun(url string, head string, token string) *CheckRun {

	type createCheckRunBody struct {
		Name    string `json:"name"`
		HeadSHA string `json:"head_sha"`
		Status  string `json:"status"`
	}

	type createCheckRunResponse struct {
		ID int `json:"id"`
	}

	requestURL := strings.ReplaceAll(url, "/hooks", "/check-runs")

	var response createCheckRunResponse

	err := RESTRequest(GitPostRequestOptions{Method: http.MethodPost, URL: requestURL, Body: createCheckRunBody{
		Name:    "Bobby",
		HeadSHA: head,
		Status:  "in_progress",
	}, Token: token}, &response)

	if err != nil {
		return nil
	}

	resUrl := fmt.Sprintf("%s/%d", requestURL, response.ID)

	checkrun := &CheckRun{URL: resUrl, Token: token}

	return checkrun
}

func (c *CheckRun) Update(status string, conclusion string) {
	type updateCheckRunBody struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}

	type updateCheckRunResponse struct {
		id string
	}

	println(status)
	var response updateCheckRunResponse
	_ = RESTRequest(GitPostRequestOptions{Method: http.MethodPatch, URL: c.URL, Body: updateCheckRunBody{
		Status:     status,
		Conclusion: conclusion,
	}, Token: c.Token}, response)

}
