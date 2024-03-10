package hostAPI

import (
	"Bobby/pkg/comm"
	"github.com/gorilla/websocket"
)

// CheckRunAPI extends HostAPI
type CheckRunAPI struct {
	Url string
	HostAPI
}

func NewCheckRunAPI(
	conn *websocket.Conn,
	body CheckRunBody,
	hooksUrl string,
	headID string,
) *CheckRunAPI {

	type NewCheckRunPayload struct {
		Body     CheckRunBody
		HooksUrl string
		HeadID   string
	}

	hostCommand := comm.HostCommand[NewCheckRunPayload]{
		Instruction: "new-checkrun-api",
		Payload: NewCheckRunPayload{
			Body:     body,
			HooksUrl: hooksUrl,
			HeadID:   headID,
		},
	}

	conn.WriteMessage(1, hostCommand.Digest())

	return &CheckRunAPI{}
}

type CheckRunOutput struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type CheckRunBody struct {
	Name    string         `json:"name"`
	HeadSHA string         `json:"head_sha"`
	Status  string         `json:"status"`
	Output  CheckRunOutput `json:"output"`
}
