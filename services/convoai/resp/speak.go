package resp

// @brief SpeakResp returned by the Conversational AI engine Speak API
//
// @since v0.9.0
type SpeakResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Successful response, see SpeakSuccessResp for details
	SuccessRes SpeakSuccessResp
}

// @brief Successful response returned by the Conversational AI engine Speak API
//
// @since v0.9.0
type SpeakSuccessResp struct {
	// The channel name of the agent
	Channel string `json:"channel"`
	// Unique identifier of the agent
	AgentId string `json:"agent_id"`
	// The start timestamp of the agent
	StartTs int64 `json:"start_ts"`
}
