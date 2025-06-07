package providers

type ReasoningEffort string

const (
	REASONING_EFFORT_HIGH   ReasoningEffort = "high"
	REASONING_EFFORT_MEDIUM ReasoningEffort = "medium"
	REASONING_EFFORT_LOW    ReasoningEffort = "low"
)

// TODO Add consturctor with options pattern
type RequestParams struct {
	UseHistory        bool
	ParallelToolCalls bool
	MaxIterations     int
	MaxTokens         int64
	Temperature       float64
	Reasoning         bool
	ReasoningEffort   ReasoningEffort
}

func NewRequestParams(options ...func(*RequestParams)) *RequestParams {
	req := &RequestParams{
		UseHistory:        false,
		ParallelToolCalls: true,
		MaxIterations:     20,
		MaxTokens:         8196,
		Temperature:       0.7,
		Reasoning:         true,
		ReasoningEffort:   REASONING_EFFORT_MEDIUM,
	}
	for _, o := range options {
		o(req)
	}
	return req
}

func WithUseHistory(enable bool) func(*RequestParams) {
	return func(req *RequestParams) {
		req.UseHistory = enable
	}
}

func WithParallelToolCalls(enable bool) func(*RequestParams) {
	return func(req *RequestParams) {
		req.ParallelToolCalls = enable
	}
}

func WithMaxIterations(iterations int) func(*RequestParams) {
	return func(req *RequestParams) {
		req.MaxIterations = iterations
	}
}

func WithMaxTokens(tokens int64) func(*RequestParams) {
	return func(req *RequestParams) {
		req.MaxTokens = tokens
	}
}

func WithTemperature(temperature float64) func(*RequestParams) {
	return func(req *RequestParams) {
		req.Temperature = temperature
	}
}

func WithReasoning(reasoning bool) func(*RequestParams) {
	return func(req *RequestParams) {
		req.Reasoning = reasoning
	}
}

func WithReasoningEffort(reasoning ReasoningEffort) func(*RequestParams) {
	return func(req *RequestParams) {
		req.ReasoningEffort = reasoning
	}
}
