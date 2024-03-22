package hostAPI

import (
	"Bobby/pkg/comm"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

// CheckRunAPI extends HostAPI
type CheckRunAPI struct {
	Url       string
	RepoId    int64
	InstallId int
	HostAPI
}

func NewCheckRunAPI(
	conn *websocket.Conn,
	body comm.CheckRunBody,
	hooksUrl string,
	repoId int64,
	installId int,
) (*CheckRunAPI, error) {

	hostCommand := comm.HostCommand[comm.NewCheckRunCommandPayload]{
		Instruction: "new-checkrun-api",
		Payload: comm.NewCheckRunCommandPayload{
			Body:         body,
			HooksUrl:     hooksUrl,
			RepositoryId: repoId,
			InstallId:    installId,
		},
	}

	hostApi := HostAPI{conn: conn, timeout: time.Minute}

	response, err := hostApi.SendCommand(&hostCommand)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Url string `json:"url"`
	}

	var checkRunPayload Response
	err = json.Unmarshal(response.Data, &checkRunPayload)

	if err != nil {
		return nil, err
	}

	if checkRunPayload.Url == "" {
		return nil, errors.New("invalid payload url")
	}

	return &CheckRunAPI{
		Url:       checkRunPayload.Url,
		RepoId:    repoId,
		InstallId: installId,
		HostAPI:   hostApi,
	}, nil
}

func (c CheckRunAPI) Update(
	status string,
	conclusion comm.CheckRunConclusion,
	output comm.CheckRunOutput,
) (bool, error) {

	hostCommand := comm.HostCommand[comm.UpdateCheckRunPayload]{
		Instruction: "update-checkrun-api",
		Payload: comm.UpdateCheckRunPayload{
			Status:       status,
			Conclusion:   conclusion,
			Output:       output,
			RepositoryId: c.RepoId,
			InstallId:    c.InstallId,
			Url:          c.Url,
		},
	}

	response, err := c.SendCommand(&hostCommand)
	type Response struct {
		Ok bool `json:"ok"`
	}

	var checkRunPayload Response
	err = json.Unmarshal(response.Data, &checkRunPayload)

	if err != nil {
		return false, err
	}

	return checkRunPayload.Ok, nil

}
