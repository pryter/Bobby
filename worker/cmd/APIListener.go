package cmd

import (
	"bobby-worker/internal/app"
	"bobby-worker/internal/events/pushEvent"
	"encoding/json"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/signal"
	"path"
)

func StartListening(config app.HTTPServiceConfig) bool {

	data, err := os.ReadFile(path.Join(config.RuntimeBasePath, "network-credentials.json"))

	if err != nil {
		ok := app.NetworkSetup(config)
		if ok {
			StartListening(config)
			return false
		}
		return ok
	}

	var networkCreds app.NetworkCredential
	err = json.Unmarshal(data, &networkCreds)

	if err != nil {
		log.Error().Msg("Unable to parse config")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme:   "ws",
		Host:     networkCreds.HostName,
		Path:     "/worker",
		RawQuery: "sid=" + networkCreds.SetupID,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to reach main service.")
		return true
	}

	log.Info().Msg("Connection created")

	defer conn.Close()

	for {
		id, message, err := conn.ReadMessage()
		if err != nil {
			return true
		}

		type WorkerPayload struct {
			SetupId string             `json:"setupId"`
			Data    github.PushPayload `json:"data"`
		}

		var payload WorkerPayload

		if id != 1 {
			continue
		}

		err = json.Unmarshal(message, &payload)

		if err != nil {
			continue
		}

		println("incoming")

		var workerEventPath = path.Join(config.RuntimeBasePath, "locker")
		pushEvent.WebhookPushEvent(
			payload.Data,
			pushEvent.WebhookPushEventOptions{RuntimeBasePath: workerEventPath},
		)
	}

}
