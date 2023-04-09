package messages

import "encoding/json"

type Action uint32

const (
	RESPONSE Action = 0
	ADD      Action = 1
	SUB      Action = 2
)

type Message struct {
	Action  Action
	Payload []int64
}

func (m *Message) Pack() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Unpack(b []byte) error {
	return json.Unmarshal(b, m)
}
