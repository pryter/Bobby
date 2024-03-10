package comm

type WorkerPayload struct {
	SetupId string      `json:"setupId"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
}
