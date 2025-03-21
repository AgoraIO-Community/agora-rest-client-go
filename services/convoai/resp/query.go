package resp

// @brief Successful response returned by the Conversational AI engine Query API
//
// @since v0.7.0
type QuerySuccessResp struct {
	// Request information
	Message string `json:"message"`
	// The agent creation timestamp
	StartTs int64 `json:"start_ts"`
	// The agent stop timestamp
	StopTs int64 `json:"stop_ts"`
	// Unique identifier of the agent
	AgentId string `json:"agent_id"`
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
	Status string `json:"status"`
}

// @brief QueryResp returned by the Conversational AI engine Query API
//
// @since v0.7.0
type QueryResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Success response, see QuerySuccessResp for details
	SuccessRes QuerySuccessResp
}
