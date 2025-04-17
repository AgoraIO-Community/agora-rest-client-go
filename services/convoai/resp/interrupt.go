package resp

// @brief InterruptResp returned by the Conversational AI engine Interrupt API
//
// @since v0.9.0
type InterruptResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Successful response, see InterruptSuccessResp for details
	SuccessRes InterruptSuccessResp
}

// @brief Successful response returned by the Conversational AI engine Interrupt API
//
// @since v0.9.0
type InterruptSuccessResp struct {
	// The channel name of the agent
	Channel string `json:"channel"`
	// Unique identifier of the agent
	AgentId string `json:"agent_id"`
	// The start timestamp of the agent
	StartTs int64 `json:"start_ts"`
}
