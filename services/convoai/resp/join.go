package resp

// @brief Successful response returned by the Conversational AI engine Join API
//
// @since v0.7.0
type JoinSuccessResp struct {
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
	Status string `json:"status"`
}

// @brief JoinResp returned by the Conversational AI engine Join API
//
// @since v0.7.0
type JoinResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Successful response, see JoinSuccessResp for details
	SuccessResp JoinSuccessResp
}
