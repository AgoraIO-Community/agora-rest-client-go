package req

// @brief Defines the request body for the Speak API
//
// @since v0.9.0
type SpeakBody struct {
	// The text content to be spoken, with a maximum length of 512 bytes.(Required)
	Text *string `json:"text"`
	// The priority of the speech behavior, which supports the following values(Optional):
	//
	//  - "INTERRUPT" (default): High priority, interrupt and speak. The agent will terminate the current interaction and speak the message directly.
	//  - "APPEND": Middle priority, append and speak. The agent will speak the message after the current interaction.
	//  - "IGNORE": Low priority, speak when idle. If the agent is currently interacting, it will directly ignore and discard the message to be spoken; only when the agent is not interacting will it speak the message.
	Priority *string `json:"priority"`
	// Whether to allow the user to speak to interrupt the agent's speech(Optional):
	//
	//  - true (default): Allow.
	//  - false: Disallow.
	Interruptable *bool `json:"interrupt"`
}
