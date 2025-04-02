package webrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireWebRecodingClientRequest struct {
	ResourceExpiredHour int      `json:"resourceExpiredHour,omitempty"`
	ExcludeResourceIds  []string `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int      `json:"regionAffinity,omitempty"`

	// StartParameter improves availability and optimizes load balancing.
	// Note: Must match the clientRequest object in subsequent start requests.
	StartParameter *StartWebRecordingClientRequest `json:"startParameter,omitempty"`
}

type QueryWebRecordingResp struct {
	api.Response
	SuccessResponse QueryWebRecordingSuccessResp
}

type QueryWebRecordingSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryWebRecordingServerResponse
}

type StartWebRecordingClientRequest struct {
	RecordingFileConfig    *api.RecordingFileConfig
	StorageConfig          *api.StorageConfig
	ExtensionServiceConfig *api.ExtensionServiceConfig
}

type UpdateWebRecordingClientRequest struct {
	WebRecordingConfig *api.UpdateWebRecordingConfig
	RtmpPublishConfig  *api.UpdateRtmpPublishConfig
}
