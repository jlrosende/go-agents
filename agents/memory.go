package agents

import "slices"

type Memory struct {
	history []string
}

func (m *Memory) Extend(messages []string) {
	m.history = append(m.history, messages...)
}

func (m *Memory) Set(messages []string) {
	m.history = slices.Clone(messages)
}

func (m *Memory) Append(message string) {
	m.history = append(m.history, message)
}

func (m Memory) Get() []string {
	return slices.Clone(m.history)
}

func (m *Memory) Clear() {
	m.history = make([]string, 0)
}
