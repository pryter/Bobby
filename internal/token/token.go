package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type tokenReqBody struct {
	RepositoryIds int64 `json:"repository_ids"`
}

func GenerateAccessToken(installationID int, repositoryID int64) (string, error) {
	postURL := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	body := tokenReqBody{repositoryID}
	parsed, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	jwtToken, err := generateJWT()

	req, err := http.NewRequest(http.MethodPost, postURL, bytes.NewBuffer(parsed))

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	return "", nil
}
