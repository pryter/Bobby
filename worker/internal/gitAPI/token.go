package gitAPI

import (
	"bobby-worker/internal/token"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// IssueAccessToken generates access token for cli clients to perform their tasks.
func IssueAccessToken(installationID int, repositoryID int64) (string, error) {
	type tokenReqBody struct {
		RepositoryIds int64 `json:"repository_ids"`
	}

	type tokenResponse struct {
		Token  string `json:"token"`
		Expire string `json:"expires_at"`
	}

	jwtToken, err := token.GenerateJWT()

	if err != nil {
		log.Error().Msg("Unable to generate local JWT for GitHub's APIs.")
		return "", err
	}

	options := GitRequestOptions{
		Method: http.MethodPost,
		URL: fmt.Sprintf(
			"https://api.github.com/app/installations/%d/access_tokens",
			installationID,
		),
		Token: jwtToken,
		Body:  tokenReqBody{RepositoryIds: repositoryID},
	}

	var response tokenResponse

	err = RESTRequest(options, &response)

	return response.Token, err
}
