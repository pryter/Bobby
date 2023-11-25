package gitAPI

import (
	"Bobby/internal/token"
	"fmt"
	"net/http"
)

func IssueAccessToken(installationID int, repositoryID int64) (string, error) {
	type tokenReqBody struct {
		RepositoryIds int64 `json:"repository_ids"`
	}

	type tokenResponse struct {
		Token  string `json:"token"`
		Expire string `json:"expires_at"`
	}

	jwtToken, err := token.GenerateJWT()

	options := GitPostRequestOptions{
		Method: http.MethodPost,
		URL:    fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID),
		Token:  jwtToken,
		Body:   tokenReqBody{RepositoryIds: repositoryID},
	}

	var response tokenResponse

	err = RESTRequest(options, &response)

	return response.Token, err
}
