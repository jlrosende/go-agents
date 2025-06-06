You are tasked with improving a response based on expert feedback. This is iteration {iteration + 1} of the refinement process.

Your goal is to address all feedback points while maintaining accuracy and relevance to the original request.

<agent:data>
<agent:request>
{request}
</agent:request>

<agent:previous-response>
{response}
</agent:previous-response>

<agent:feedback>
<rating>{feedback.rating}</rating>

<details>{feedback.feedback}</details>
<focus-areas>{focus_areas}</focus-areas>
</agent:feedback>
</agent:data>

<agent:instruction>
Create an improved version of the response that:

1. Directly addresses each point in the feedback
2. Focuses on the specific areas mentioned for improvement
3. Maintains all the strengths of the original response
4. Remains accurate and relevant to the original request

Provide your complete improved response without explanations or commentary.
</agent:instruction>
