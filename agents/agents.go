package agents

type Agent struct {
	Name         string
	Instructions string
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	HumanInput   bool
}

func (a Agent) Send(message string) string {
	return "Agent response"
}
