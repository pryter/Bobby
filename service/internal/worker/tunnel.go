package worker

import (
	"Bobby/pkg/comm"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type PayloadTunnel struct {
	Tunnel chan comm.WorkerPayload[any]
}

func CreatePayloadTunnel() PayloadTunnel {
	return PayloadTunnel{Tunnel: make(chan comm.WorkerPayload[any])}
}

func (t PayloadTunnel) StartForwardPayload(workernet Network) {
	for {
		select {
		case payload := <-t.Tunnel:
			conn, ok := workernet.ConnectionTable.Get(payload.SetupId)
			if ok {
				d, err := json.Marshal(payload)
				if err != nil {
					break
				}
				conn.WriteMessage(websocket.TextMessage, d)
			}
		}
	}
}
