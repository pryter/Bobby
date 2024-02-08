package worker

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type WorkerNetwork struct {
	ConnectionTable ConnectionTable
	WSUpgrader      websocket.Upgrader
}

type HostCommand struct {
	Instruction string      `json:"instruction"`
	Payload     interface{} `json:"payload"`
}

func createPingHandler(c *websocket.Conn, id string, duration time.Duration) chan struct{} {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	handler := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.WriteMessage(websocket.PingMessage, []byte("Ping!"))
				if err != nil {
					log.Printf("Ping server destroyed with error #%s", id)
					return
				}
				break
			case <-handler:
				log.Printf("Ping server destroyed #%s", id)
				return
			}
		}
	}()

	return handler
}

func (n WorkerNetwork) onMessageReceived(id int, message []byte) {
	if id != 1 {
		return
	}

	var command HostCommand
	err := json.Unmarshal(message, &command)

	if err != nil {
		log.Error().Str("instruction", command.Instruction).Err(err).Msg("Unrecognised command.")
		return
	}

	switch command.Instruction {
	case "registere":
		break
	default:
		log.Error().Str("instruction", command.Instruction).Msg("Unrecognised command.")
	}
}

func (n WorkerNetwork) HttpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("sid")

	c, err := n.WSUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msgf("Unable to upgrade the connection #%s.", id)
		return
	}

	log.Debug().Msgf("Connection created from #%s", id)
	n.ConnectionTable.Set(id, c)

	defer close(createPingHandler(c, id, time.Second*10))

	for {
		mid, message, err := c.ReadMessage()

		if err != nil {
			log.Printf("Connection closed from #%s", id)
			n.ConnectionTable.Remove(id)
			return
		}

		n.onMessageReceived(mid, message)
	}
}
