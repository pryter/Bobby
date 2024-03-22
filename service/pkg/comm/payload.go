package comm

import (
	"encoding/json"
	"github.com/go-playground/webhooks/v6/github"
)

type WorkerPayload[T interface{}] struct {
	SetupId   string `json:"setupId"`
	Action    string `json:"action"`
	PayloadId string `json:"payloadId"`
	Data      T      `json:"data"`
}

func (w WorkerPayload[T]) Digest() ([]byte, error) {
	return json.Marshal(w)
}

/* Pre-defined Worker Payloads */

type PushEventWorkerPayload = WorkerPayload[github.PushPayload]
type SetupEventWorkerPayload = WorkerPayload[string]

/* Pre-defined HostCommand Payloads */

type RegisterCommandPayload struct {
	MacAddr string `json:"mac-addr"`
}

type CheckRunOutput struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type CheckRunBody struct {
	Name    string         `json:"name"`
	HeadSHA string         `json:"head_sha"`
	Status  string         `json:"status"`
	Output  CheckRunOutput `json:"output"`
}

type CheckRunConclusion string

var (
	ConclusionSuccess   CheckRunConclusion = "success"
	ConclusionFailure   CheckRunConclusion = "failure"
	ConclusionCancelled CheckRunConclusion = "cancelled"
	ConclusionTimedOut  CheckRunConclusion = "timed_out"
)

type UpdateCheckRunPayload struct {
	Status       string             `json:"status"`
	Conclusion   CheckRunConclusion `json:"conclusion"`
	Output       CheckRunOutput     `json:"output"`
	RepositoryId int64              `json:"repositoryId"`
	InstallId    int                `json:"installId"`
	Url          string             `json:"url"`
}

type NewCheckRunCommandPayload struct {
	RepositoryId int64        `json:"repositoryId"`
	InstallId    int          `json:"installId"`
	Body         CheckRunBody `json:"body"`
	HooksUrl     string       `json:"hooksUrl"`
}
