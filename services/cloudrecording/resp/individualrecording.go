package resp

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief QueryIndividualRecordingResp returned by the individual recording Query API.
//
// @since v0.8.0
type QueryIndividualRecordingResp struct {
	// Response returned by the cloud recording API, see Response for details
	api.Response
	// Success response, see QueryIndividualRecordingSuccessResp for details
	SuccessRes QueryIndividualRecordingSuccessResp
}

// @brief Successful response returned by the individual recording Query API.
//
// @since v0.8.0
type QueryIndividualRecordingSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string
	// Unique identifier of the recording session
	Sid string
	// Server response, see QueryIndividualRecordingServerResponse for details
	ServerResponse *api.QueryIndividualRecordingServerResponse
}

// @brief QueryIndividualRecordingVideoScreenshotResp returned by the individual recording QueryVideoScreenshot API.
//
// @since v0.8.0
type QueryIndividualRecordingVideoScreenshotResp struct {
	// Response returned by the cloud recording API, see Response for details
	api.Response
	// Success response, see QueryIndividualRecordingVideoScreenshotSuccessResp for details
	SuccessRes QueryIndividualRecordingVideoScreenshotSuccessResp
}

// @brief Successful response returned by the individual recording QueryVideoScreenshot API.
//
// @since v0.8.0
type QueryIndividualRecordingVideoScreenshotSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string
	// Unique identifier of the recording session
	Sid string
	// Server response, see QueryIndividualVideoScreenshotServerResponse for details
	ServerResponse *api.QueryIndividualVideoScreenshotServerResponse
}
