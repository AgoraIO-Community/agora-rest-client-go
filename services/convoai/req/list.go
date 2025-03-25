package req

type ListOptions struct {
	Limit    *int
	State    *int
	FromTime *int
	ToTime   *int
	Cursor   *string
	Channel  *string
}

// @brief ListOption Define the filter condition type used to query the list of intelligent agents
//
// @since v0.7.0
type ListOption func(*ListOptions)

// @brief WithLimit Set the number of items.
//
// @param Limit The number of items returned per page for pagination,default is 20
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithLimit(limit int) ListOption {
	return func(opts *ListOptions) {
		opts.Limit = &limit
	}
}

// @brief WithState Set the state.
//
// @param State Status of the intelligent agent to query，supports the following statuses:
//
//   - IDLE (0): Idle State of the intelligent agent
//
//   - STARTING (1): Intelligent agent is starting
//
//   - RUNNING (2): Intelligent agent is running
//
//   - STOPPING (3): Intelligent agent is stopping
//
//   - STOPPED (4): Intelligent agent has completed exit
//
//   - RECOVERING (5): The agent is recovering.
//
//   - FAILED (6): The agent has failed.
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithState(state int) ListOption {
	return func(opts *ListOptions) {
		opts.State = &state
	}
}

// @brief WithFromTime Set the start timestamp.
//
// @param FromTime Start timestamp (s), default is 1 day ago
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithFromTime(fromTime int) ListOption {
	return func(opts *ListOptions) {
		opts.FromTime = &fromTime
	}
}

// @brief WithToTime Set the end timestamp.
//
// @param ToTime End timestamp (s), default is the current time
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithToTime(toTime int) ListOption {
	return func(opts *ListOptions) {
		opts.ToTime = &toTime
	}
}

// @brief WithCursor Set the pagination cursor.
//
// @param Cursor Pagination Cursor, i.e., the intelligent agent ID of the pagination start position
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithCursor(cursor string) ListOption {
	return func(opts *ListOptions) {
		opts.Cursor = &cursor
	}
}

// @brief WithChannel Set the channel name.
//
// @param channel Channel name
//
// @return Returns the ListOption function
//
// @since v0.7.0
func WithChannel(channel string) ListOption {
	return func(opts *ListOptions) {
		opts.Channel = &channel
	}
}
