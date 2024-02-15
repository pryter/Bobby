package comm

type WorkerPayload struct {
	SetupId string      `json:"setupId"`
	Data    interface{} `json:"data"`
}
