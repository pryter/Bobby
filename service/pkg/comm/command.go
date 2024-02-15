package comm

import (
	"encoding/json"
)

type RegisterCommandPayload struct {
	MacAddr string `json:"mac-addr"`
}

type HostCommand[T interface{}] struct {
	Instruction string `json:"instruction"`
	Payload     T      `json:"payload"`
}

type RawHostCommand struct {
	Instruction string          `json:"instruction"`
	Payload     json.RawMessage `json:"payload"`
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

func (c HostCommand[T]) Digest() []byte {
	bytes, err := json.Marshal(
		map[string]interface{}{
			"instruction": c.Instruction,
			"payload":     c.Payload,
		},
	)

	if err != nil {
		panic(err)
	}

	return bytes
}
