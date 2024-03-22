package comm

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Digestible interface {
	Digest() []byte
	GetId() string
}

type HostCommand[T interface{}] struct {
	Instruction   string `json:"instruction"`
	Payload       T      `json:"payload"`
	TransactionId string `json:"transactionId"`
}

type RawHostCommand struct {
	Instruction   string          `json:"instruction"`
	Payload       json.RawMessage `json:"payload"`
	TransactionId string          `json:"transactionId"`
}

func (c RawHostCommand) ResolvePayload(v any) error {

	err := json.Unmarshal(c.Payload, v)

	if err != nil {
		return err
	}

	return nil
}

func ParseHostCommand(payload []byte) (RawHostCommand, error) {
	var command RawHostCommand
	err := json.Unmarshal(payload, &command)
	if err != nil {
		return RawHostCommand{}, err
	}

	return command, nil
}

func (c *HostCommand[T]) GetId() string {
	return c.TransactionId
}

func (c *HostCommand[T]) Digest() []byte {

	id, err := uuid.NewUUID()
	c.TransactionId = id.String()

	bytes, err := json.Marshal(
		map[string]interface{}{
			"instruction":   c.Instruction,
			"payload":       c.Payload,
			"transactionId": c.TransactionId,
		},
	)

	if err != nil {
		panic(err)
	}

	return bytes
}
