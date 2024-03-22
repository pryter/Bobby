package hostAPI

import (
	"Bobby/pkg/comm"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

// HostAPI are REST API
type HostAPI struct {
	conn    *websocket.Conn
	timeout time.Duration
}

func (t *HostAPI) SendCommand(command comm.Digestible) (comm.WorkerPayload[json.RawMessage], error) {
	digested := command.Digest()
	err := t.conn.WriteMessage(1, digested)
	if err != nil {
		return comm.WorkerPayload[json.RawMessage]{}, err
	}

	res, err := t.awaitForResponse(command.GetId())
	if err != nil {
		return comm.WorkerPayload[json.RawMessage]{}, err
	}

	return res, nil
}

func (t *HostAPI) awaitForResponse(
	tid string,
) (comm.WorkerPayload[json.RawMessage], error) {

L:
	for timeout := time.After(t.timeout); ; {
		select {
		case <-timeout:
			break L
		default:
			{
				id, message, err := t.conn.ReadMessage()
				if err != nil {
					return comm.WorkerPayload[json.RawMessage]{}, err
				}

				var payload comm.WorkerPayload[json.RawMessage]

				if id != 1 {
					continue
				}

				err = json.Unmarshal(message, &payload)

				if err != nil {
					continue
				}

				// matched response
				if payload.PayloadId == tid {
					return payload, nil
				}
			}
		}
	}

	return comm.WorkerPayload[json.RawMessage]{}, errors.New("timeout")
}
