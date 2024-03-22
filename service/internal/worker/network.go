package worker

import (
	"Bobby/internal/app"
	"Bobby/pkg/challenge"
	"Bobby/pkg/comm"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"path"
	"time"
)

type Network struct {
	ConnectionTable ConnectionTable
	WSUpgrader      websocket.Upgrader
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

func closeWithMessage(c *websocket.Conn, message string) {
	cm := websocket.FormatCloseMessage(
		1000, message,
	)
	c.WriteMessage(8, cm)
	c.Close()
}

func (n Network) onMessageReceived(setupId string, id int, message []byte, conn *websocket.Conn) {
	if id != 1 {
		return
	}

	command, err := comm.ParseHostCommand(message)

	if err != nil {
		log.Error().Str("instruction", command.Instruction).Err(err).Msg("Unrecognised command.")
		return
	}

	switch command.Instruction {
	case "register":
		registerAction(conn, command)
		break
	case "new-checkrun-api":
		newCheckRunApiAction(conn, setupId, command)
		break
	case "update-checkrun-api":
		updateCheckRunApiAction(conn, setupId, command)
		break
	default:
		log.Error().Str("instruction", command.Instruction).Msg("Unrecognised command.")
	}
}

func (n Network) HttpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("sid")

	c, err := n.WSUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msgf("Unable to upgrade the connection #%s.", id)
		return
	}

	// For registered clients
	if id != "" {

		challengeString := r.URL.Query().Get("challenge")

		if challengeString == "" {
			closeWithMessage(c, "the connection must contains a valid challenge")
			return
		}

		// Solve the connection challenge
		workerDataPath := app.GetResources().WorkerData.GetAbsolutePath()
		workerFolder := path.Join(workerDataPath, id)
		challengeSolver := challenge.Solver{
			SetupID:        id,
			PrivateKeyPath: path.Join(workerFolder, "key"),
			SecretPath:     path.Join(workerFolder, "secret"),
		}

		ok, err := challengeSolver.Solve(challengeString)

		// Reject connection if received invalid challenge
		if !ok || err != nil {
			closeWithMessage(c, "unable to solve the given challenge")
			return
		}

		log.Debug().Msgf("Connection created from #%s", id)
		n.ConnectionTable.Set(id, c)

		defer close(createPingHandler(c, id, time.Second*10))
	}

	for {
		// Set timeout for
		if id == "" {
			time.AfterFunc(
				time.Minute*2, func() {
					closeWithMessage(c, "connection timeout")
				},
			)
		}

		mid, message, err := c.ReadMessage()

		if err != nil {
			// Ignore connection err from unregistered connections
			if id == "" {
				return
			}

			log.Printf("Connection closed from #%s", id)
			n.ConnectionTable.Remove(id)
			return
		}

		n.onMessageReceived(id, mid, message, c)
	}
}
