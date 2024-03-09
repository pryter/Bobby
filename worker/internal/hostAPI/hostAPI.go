package hostAPI

import (
	"bobby-worker/internal/events/pushEvent"
	"encoding/json"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/gorilla/websocket"
)

// HostAPI are REST API
type HostAPI struct {
	conn *websocket.Conn
}

func (t *HostAPI) awaitForResponse(key string) {
	for {
		id, message, err := t.conn.ReadMessage()
		if err != nil {
			return
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

		var workerEventPath = ""
		pushEvent.WebhookPushEvent(
			payload.Data,
			pushEvent.WebhookPushEventOptions{RuntimeBasePath: workerEventPath},
		)
	}
}
