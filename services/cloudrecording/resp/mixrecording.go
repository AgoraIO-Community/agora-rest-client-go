package resp

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief Successful response returned by the mix recording QueryHLS API.
//
// @since v0.8.0
type QueryMixRecordingHLSSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string
	// Unique identifier of the recording session
	Sid string
	// Server response, see QueryMixRecordingHLSServerResponse for details
	ServerResponse *api.QueryMixRecordingHLSServerResponse
}

// @brief Response returned by the mix recording QueryHLS API.
//
// @since v0.8.0
type QueryMixRecordingHLSResp struct {
	// Response returned by the cloud recording API, see Response for details
	api.Response
	// Successful response, see QueryMixRecordingHLSSuccessResp for details
	SuccessResponse QueryMixRecordingHLSSuccessResp
}

// @brief Successful response returned by the mix recording QueryHLSAndMP4 API.
//
// @since v0.8.0
type QueryMixRecordingHLSAndMP4SuccessResp struct {
	// Unique identifier of the resource
	ResourceId string
	// Unique identifier of the recording session
	Sid string
	// Server response, see QueryMixRecordingHLSAndMP4ServerResponse for details
	ServerResponse *api.QueryMixRecordingHLSAndMP4ServerResponse
}

// @brief Response returned by the mix recording QueryHLSAndMP4 API.
//
// @since v0.8.0
type QueryMixRecordingHLSAndMP4Resp struct {
	// Response returned by the cloud recording API, see Response for details
	api.Response
	// Successful response, see QueryMixRecordingHLSAndMP4SuccessResp for details
	SuccessResponse QueryMixRecordingHLSAndMP4SuccessResp
}
