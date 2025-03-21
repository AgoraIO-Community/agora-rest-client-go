package resp

// @brief Successful response returned by the Conversational AI engine List API
//
// @since v0.7.0
type ListSuccessResp struct {
	// Intelligent agent data information
	Data struct {
		// Number of intelligent agents returned this time
		Count int `json:"count"`
		// List of intelligent agents that meet the conditions
		List []struct {
			// Intelligent agent creation timestamp
			StartTs int64 `json:"start_ts"`
			// Intelligent agent running status
			//
			//  - IDLE (0): Idle state of the intelligent agent.
			//
			//  - STARTING (1): Intelligent agent is starting.
			//
			//  - RUNNING (2): Intelligent agent is running.
			//
			//  - STOPPING (3): Intelligent agent is stopping.
			//
			//  - STOPPED (4): Intelligent agent has completed exit.
			//
			//  - RECOVERING (5): Intelligent agent is recovering.
			//
			//  - FAILED (6): Intelligent agent execution failed.
			Status string `json:"status"`
			// Unique identifier of the intelligent agent
			AgentId string `json:"agent_id"`
		} `json:"list"`
	} `json:"data"`
	// Metadata of the returned list
	Meta struct {
		// Pagination cursor
		Cursor string `json:"cursor"`
		// Total number of intelligent agents that meet the query conditions
		Total int `json:"total"`
	} `json:"meta"`
	// Request status
	Status string `json:"status"`
}

// @brief ListResp returned by the Conversational AI engine List API
//
// @since v0.7.0
type ListResp struct {
	// Response returned by the Conversational AI engine API, see Response for details
	Response
	// Successful response, see ListSuccessResp for details
	SuccessRes ListSuccessResp
}
