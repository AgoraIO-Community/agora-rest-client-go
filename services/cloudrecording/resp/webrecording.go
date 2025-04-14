package resp

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief Response returned by the web recording Query API.
//
// @since v0.8.0
type QueryWebRecordingResp struct {
	// Response returned by the cloud recording API, see Response for details
	api.Response
	// Success response, see QueryWebRecordingSuccessResp for details
	SuccessResponse QueryWebRecordingSuccessResp
}

// @brief Successful response returned by the web recording Query API.
//
// @since v0.8.0
type QueryWebRecordingSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string
	// Unique identifier of the recording session
	Sid string
	// Server response, see QueryWebRecordingServerResponse for details
	ServerResponse *api.QueryWebRecordingServerResponse
}
