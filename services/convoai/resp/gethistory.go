package resp

// @brief HistoryResp returned by the Conversational AI engine History API
//
// @since v0.9.0
type HistoryResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Success response, see HistorySuccessResp for details
	SuccessRes HistorySuccessResp
}

// @brief Successful response returned by the Conversational AI engine History API
//
// @since v0.9.0
type HistorySuccessResp struct {
	// The agent creation timestamp
	StartTs int64 `json:"start_ts"`
	// Unique identifier of the agent
	AgentId string `json:"agent_id"`
	// Only returns the running status of the agent
	Status string `json:"status"`
	// Agent short-term memory content
	Contents []HistoryContent `json:"contents"`
}

// @brief Agent short-term memory content
//
// @since v0.9.0
type HistoryContent struct {
	// The role of sending messages:
	//
	//  - user: User.
	//
	//  - assistant: Agent.
	Role string `json:"role"`
	// The content of the message
	Content string `json:"content"`
}
