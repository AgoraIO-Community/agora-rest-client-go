package resp

// @brief Successful response returned by the Conversational AI engine Update API
//
// @since v0.7.0
type UpdateSuccessResp struct {
	// Unique identifier of the agent
	AgentId string `json:"agent_id"`
	// Timestamp when the agent was created
	CreateTs int `json:"create_ts"`
	// Running status of the agent
	//
	//  - IDLE (0): The agent is idle.
	//
	//  - STARTING (1): The agent is starting.
	//
	//  - RUNNING (2): The agent is running.
	//
	//  - STOPPING (3): The agent is stopping.
	//
	//  - STOPPED (4): The agent has stopped.
	//
	//  - RECOVERING (5): The agent is recovering.
	//
	//  - FAILED (6): The agent has failed.
	State string `json:"state"`
}

// @brief Response returned by the Conversational AI engine Update API
//
// @since v0.7.0
type UpdateResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Successful response, see UpdateSuccessResp for details
	SuccessResp UpdateSuccessResp
}
