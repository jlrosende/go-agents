package router

type Plan struct {
	Steps      []Step `json:"steps" jsonschema_description:"List of steps to execute sequentially"`
	IsComplete bool   `json:"is_complete" jsonschema_description:"Whether the overall plan objective is complete"`
}

type Step struct {
	Description string `json:"description" jsonschema_description:"Description of the step"`
	Tasks       []Task `json:"tasks" jsonschema_description:"Subtasks that can be executed in parallel"`
}

type Task struct {
	Description string `json:"description" jsonschema_description:"Subtasks that can be executed in parallel"`
	Agent       string `json:"agent" jsonschema_description:"Name of Agent from given list of agents that the LLM has access to for this task"`
}
