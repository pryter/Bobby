package worker

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type WorkerPayload struct {
	SetupId string      `json:"setupId"`
	Data    interface{} `json:"data"`
}

type PayloadTunnel struct {
	Tunnel chan WorkerPayload
}

func CreatePayloadTunnel() PayloadTunnel {
	return PayloadTunnel{Tunnel: make(chan WorkerPayload)}
}

func (t PayloadTunnel) StartForwardPayload(workernet WorkerNetwork) {
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
