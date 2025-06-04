package req

// @brief Request body for calling the Conversational AI engine Update API
//
// @since v0.7.0
type UpdateReqBody struct {
	// Dynamic key (Token) for authentication. If your project has enabled App Certificate, you must pass the dynamic key of your project in this field.
	Token string `json:"token"`

	// LLM configuration for the agent
	LLM *UpdateLLMBody `json:"llm,omitempty"`
}

// @brief Defines the LLM configuration for the agent
//
// @since v0.11.0
type UpdateLLMBody struct {
	// A set of predefined information attached at the beginning of each LLM call to control LLM output (optional)
	//
	// Can be role settings, prompts, and answer samples, must be compatible with the OpenAI protocol
	SystemMessages []map[string]any `json:"system_messages,omitempty"`

	// Additional information transmitted in the LLM message body, such as the model used, maximum token Limit, etc. (optional)
	//
	// Different LLM vendors support different configurations, see their respective LLM documentation for details
	Params map[string]any `json:"params,omitempty"`
}
