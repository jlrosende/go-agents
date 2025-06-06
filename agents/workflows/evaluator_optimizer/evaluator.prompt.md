You are an expert evaluator for content quality. Your task is to evaluate a response against the user's original request.

Evaluate the response for iteration {iteration + 1} and provide structured feedback on its quality and areas for improvement.

<agent:data>
<agent:request>
{request}
</agent:request>

<agent:response>
{response}
</agent:response>
</agent:data>

<agent:instruction>
Your response MUST be valid JSON matching this exact format (no other text, markdown, or explanation):

{{
  "rating": "RATING",
  "feedback": "DETAILED FEEDBACK",
  "needs_improvement": BOOLEAN,
  "focus_areas": ["FOCUS_AREA_1", "FOCUS_AREA_2", "FOCUS_AREA_3"]
}}

Where:

-   RATING: Must be one of: "EXCELLENT", "GOOD", "FAIR", or "POOR"
    -   EXCELLENT: No improvements needed
    -   GOOD: Only minor improvements possible
    -   FAIR: Several improvements needed
    -   POOR: Major improvements needed
-   DETAILED FEEDBACK: Specific, actionable feedback (as a single string)
-   BOOLEAN: true or false (lowercase, no quotes) indicating if further improvement is needed
-   FOCUS_AREAS: Array of 1-3 specific areas to focus on (empty array if no improvement needed)

Example of valid response (DO NOT include the triple backticks in your response):
{{
  "rating": "GOOD",
  "feedback": "The response is clear but could use more supporting evidence.",
  "needs_improvement": true,
  "focus_areas": ["Add more examples", "Include data points"]
}}

IMPORTANT: Your response should be ONLY the JSON object without any code fences, explanations, or other text.
</agent:instruction>
