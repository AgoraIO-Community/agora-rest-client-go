package individualrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireIndividualRecodingClientRequest struct {
	// Resource expiration time in hours. Default: 72 hours.
	ResourceExpiredHour int

	// Resource IDs to exclude, enabling cross-region recording.
	ExcludeResourceIds []string

	// Region affinity for recording resources:
	// 0: Closest to request origin
	// 1: China
	// 2: Southeast Asia
	// 3: Europe
	// 4: North America
	RegionAffinity int

	// StartParameter improves availability and optimizes load balancing.
	// Note: Must match the clientRequest object in subsequent start requests.
	StartParameter *StartIndividualRecordingClientRequest
}

type StartIndividualRecordingClientRequest struct {
	// Authentication token
	Token string

	// Third-party cloud storage configuration
	StorageConfig *api.StorageConfig

	// Audio/video stream recording configuration
	RecordingConfig *api.RecordingConfig

	// Recording file configuration
	RecordingFileConfig *api.RecordingFileConfig

	// Video snapshot configuration
	SnapshotConfig *api.SnapshotConfig

	// Application configuration
	AppsCollection *api.AppsCollection

	// Transcoding options for delayed transcoding or mixing
	TranscodeOptions *api.TranscodeOptions
}

type QueryIndividualRecordingSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryIndividualRecordingServerResponse
}

type QueryIndividualRecordingResp struct {
	api.Response
	SuccessResponse QueryIndividualRecordingSuccessResp
}

type QueryIndividualRecordingVideoScreenshotSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryIndividualVideoScreenshotServerResponse
}

type QueryIndividualRecordingVideoScreenshotResp struct {
	api.Response
	SuccessResponse QueryIndividualRecordingVideoScreenshotSuccessResp
}

type UpdateIndividualRecordingClientRequest struct {
	StreamSubscribe *api.UpdateStreamSubscribe
}
