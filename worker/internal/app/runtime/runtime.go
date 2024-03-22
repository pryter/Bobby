package runtime

import (
	"Bobby/pkg/comm"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/signal"
)

type WorkerRuntime struct {
	ConnectionUrl url.URL
	events        map[string]func(rawPayload json.RawMessage, conn *websocket.Conn)
}

func (r *WorkerRuntime) connect() (bool, *websocket.Conn) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(r.ConnectionUrl.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to reach main service.")
		return false, nil
	}

	log.Info().Msg("Connection created")
	return true, conn
}

func (r *WorkerRuntime) waitForPayload(conn *websocket.Conn) bool {
	for {
		id, message, err := conn.ReadMessage()
		if err != nil {
			return false
		}

		var payload comm.WorkerPayload[json.RawMessage]

		if id != 1 {
			continue
		}

		err = json.Unmarshal(message, &payload)

		// Skip payload if cannot parse
		if err != nil {
			continue
		}

		e, ok := r.events[payload.Action]
		if !ok {
			log.Error().Msg("No registered event named " + payload.Action)
		}

		e(payload.Data, conn)
	}
}

func (r *WorkerRuntime) Run() bool {
	ok, conn := r.connect()
	if !ok {
		return false
	}
	defer conn.Close()

	r.waitForPayload(conn)
	return true
}

func (r *WorkerRuntime) RegisterEvent(
	name string,
	event func(rawPayload json.RawMessage, conn *websocket.Conn),
) {
	if r.events == nil {
		r.events = make(map[string]func(rawPayload json.RawMessage, conn *websocket.Conn))
	}

	r.events[name] = event
}
