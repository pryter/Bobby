package gitAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitPostRequestOptions struct {
	Method string
	Body   interface{}
	URL    string
	Token  string
}

func RESTRequest(options GitPostRequestOptions, response any) error {
	parsed, err := json.Marshal(options.Body)
	req, err := http.NewRequest(options.Method, options.URL, bytes.NewBuffer(parsed))

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", options.Token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(response)
	return err
}
