package worker

import (
	"Bobby/pkg/comm"
	"Bobby/pkg/gitAPI"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func newCheckRunApiAction(conn *websocket.Conn, setupId string, command comm.RawHostCommand) {
	respond := func(url string) {
		payload := comm.WorkerPayload[struct {
			Url string `json:"url"`
		}]{
			SetupId:   setupId,
			Action:    "response-api",
			PayloadId: command.TransactionId,
			Data: struct {
				Url string `json:"url"`
			}{Url: url},
		}

		digested, err := payload.Digest()
		if err != nil {
			log.Debug().Err(err).Msg("unable to digest respond payload")
			return
		}

		conn.WriteMessage(1, digested)
	}

	var cp comm.NewCheckRunCommandPayload
	err := command.ResolvePayload(&cp)

	if err != nil {
		respond("")
		return
	}

	accessToken, _ := gitAPI.IssueAccessToken(cp.InstallId, cp.RepositoryId)

	url, err := gitAPI.NewCheckRun(cp.HooksUrl, cp.Body, accessToken)

	if err != nil {
		respond("")
		return
	}

	respond(url)
}

func updateCheckRunApiAction(conn *websocket.Conn, setupId string, command comm.RawHostCommand) {
	respond := func(ok bool) {
		payload := comm.WorkerPayload[struct {
			Ok bool `json:"ok"`
		}]{
			SetupId:   setupId,
			Action:    "response-api",
			PayloadId: command.TransactionId,
			Data: struct {
				Ok bool `json:"ok"`
			}{Ok: ok},
		}

		digested, err := payload.Digest()
		if err != nil {
			log.Debug().Err(err).Msg("unable to digest respond payload")
			return
		}

		conn.WriteMessage(1, digested)
	}

	var cp comm.UpdateCheckRunPayload
	err := command.ResolvePayload(&cp)

	if err != nil {
		respond(false)
		return
	}

	accessToken, _ := gitAPI.IssueAccessToken(cp.InstallId, cp.RepositoryId)
	gitAPI.UpdateCheckRun(cp, accessToken)

	respond(true)
}
