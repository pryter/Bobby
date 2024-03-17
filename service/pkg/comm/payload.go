package comm

import (
	"encoding/json"
	"github.com/go-playground/webhooks/v6/github"
)

type WorkerPayload[T interface{}] struct {
	SetupId string `json:"setupId"`
	Action  string `json:"action"`
	Data    T      `json:"data"`
}

func (w WorkerPayload[T]) Digest() ([]byte, error) {
	return json.Marshal(w)
}

/* Pre-defined Worker Payloads */

type PushEventWorkerPayload = WorkerPayload[github.PushPayload]
type SetupEventWorkerPayload = WorkerPayload[string]
