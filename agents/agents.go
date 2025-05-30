package agents

type Agent struct {
	Name         string
	Instructions string
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	// HumanInput   bool
	Model string
}

func (a Agent) Send(message string) string {
	return "Agent response"
}

func (a Agent) Generate(message string) string {
	return "Agent response"
}

func (a Agent) Strcutured(message string) string {
	return "Agent response"
}
