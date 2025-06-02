package memory

import "slices"

type Memory struct {
	history []Message
}

func (m *Memory) Extend(messages []Message) {
	m.history = append(m.history, messages...)
}

func (m *Memory) Set(messages []Message) {
	m.history = slices.Clone(messages)
}

func (m *Memory) Append(message Message) {
	m.history = append(m.history, message)
}

func (m Memory) Get() []Message {
	return slices.Clone(m.history)
}

func (m *Memory) Clear() {
	m.history = make([]Message, 0)
}

type MessageType string

const (
	MESSAGE_TYPE_DEVELOPER MessageType = "developer"
	MESSAGE_TYPE_SYSTEM    MessageType = "system"
	MESSAGE_TYPE_USER      MessageType = "user"
	MESSAGE_TYPE_ASSISTANT MessageType = "assistant"
	MESSAGE_TYPE_TOOL      MessageType = "tool"
)

type Message struct {
	Type       MessageType
	Content    string
	ToolCallID string
}
