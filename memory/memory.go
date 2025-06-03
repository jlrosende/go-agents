package memory

import "slices"

type Memory struct {
	history []any
}

func (m *Memory) Extend(messages []any) {
	m.history = append(m.history, messages...)
}

func (m *Memory) Set(messages []any) {
	m.history = slices.Clone(messages)
}

func (m *Memory) Append(message any) {
	m.history = append(m.history, message)
}

func (m Memory) Get() []any {
	return slices.Clone(m.history)
}

func (m *Memory) Clear() {
	m.history = make([]any, 0)
}
