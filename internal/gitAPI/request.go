package gitAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// GitRequestOptions is options for created git api request.
type GitRequestOptions struct {
	Method string
	Body   interface{}
	URL    string
	Token  string
}

func RESTRequest(options GitRequestOptions, response any) error {
	parsed, err := json.Marshal(options.Body)
	req, err := http.NewRequest(
		options.Method, options.URL, bytes.NewBuffer(parsed),
	)

	// specified required headers for GitHub's api requests.
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", options.Token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal().Err(err).Str(
			"path", options.URL,
		).Str("method", options.Method).Msg("Unable to request git api.")
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Unable to close the http request.")
		}
	}()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		log.Fatal().Fields(
			map[string]string{
				"code":   strconv.Itoa(res.StatusCode),
				"status": res.Status,
			},
		).Msg("Unusable response.")
		return err
	}

	err = json.NewDecoder(res.Body).Decode(response)
	return err
}
