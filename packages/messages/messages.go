package messages

import "encoding/json"

type Action uint8

const (
	ERROR    Action = 0
	ADD      Action = 1
	SUB      Action = 2
	RESPONSE Action = 3
	LOCK     Action = 4
	UNLOCK   Action = 5
	INUSE    Action = 6
	GRANTED  Action = 7
)

type Message struct {
	Action  Action
	Payload map[string]interface{}
}

func (m *Message) Pack() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Unpack(b []byte) error {
	return json.Unmarshal(b, m)
}
