package evaluator_optimizer

type Rating string

const (
	RATING_POOR     Rating = "poor"
	RATING_FAIR     Rating = "fair"
	RATING_GOOD     Rating = "good"
	RATING_EXCELENT Rating = "execelent"
)

type Evaluation struct {
	Rating           Rating   `json:"rating"  jsonschema:"enum=poor,enum=fair,enum=good,enum=excelent" jsonschema_description:"Quality rating of the response"`
	Feedback         string   `json:"feedback" jsonschema_description:"Specific feedback and suggestions for improvement"`
	NeedsImprovement bool     `json:"needs_improvement" jsonschema_description:"Whether the output needs further improvement"`
	FocusAreas       []string `json:"focus_areas" jsonschema_description:"Specific areas to focus on in next iteration"`
}
