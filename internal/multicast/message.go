package multicast

import (
	"encoding/json"
	"time"
)

// Message represents a multicast test message
type Message struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// Marshal serializes the message to JSON
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalMessage deserializes JSON data into a Message
func UnmarshalMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// Age returns how long ago the message was created
func (m *Message) Age() time.Duration {
	return time.Since(m.Timestamp)
}
