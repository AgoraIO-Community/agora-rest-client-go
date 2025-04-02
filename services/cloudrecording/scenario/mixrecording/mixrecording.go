package mixrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireMixRecodingClientRequest struct {
	ResourceExpiredHour int
	ExcludeResourceIds  []string
	RegionAffinity      int

	// StartParameter improves availability and optimizes load balancing.
	// Note: Must match the clientRequest object in subsequent start requests.
	StartParameter *StartMixRecordingClientRequest
}

type QueryMixRecordingHLSSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryMixRecordingHLSServerResponse
}

type QueryMixRecordingHLSResp struct {
	api.Response
	SuccessResponse QueryMixRecordingHLSSuccessResp
}

type QueryMixRecordingHLSAndMP4SuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryMixRecordingHLSAndMP4ServerResponse
}
type QueryMixRecordingHLSAndMP4Resp struct {
	api.Response
	SuccessResponse QueryMixRecordingHLSAndMP4SuccessResp
}

type StartMixRecordingClientRequest struct {
	Token               string
	RecordingConfig     *api.RecordingConfig
	RecordingFileConfig *api.RecordingFileConfig
	StorageConfig       *api.StorageConfig
}

type UpdateMixRecordingClientRequest struct {
	StreamSubscribe *api.UpdateStreamSubscribe
}

type UpdateLayoutUpdateMixRecordingClientRequest struct {
	MaxResolutionUID           string
	MixedVideoLayout           int
	BackgroundColor            string
	BackgroundImage            string
	DefaultUserBackgroundImage string
	LayoutConfig               []api.UpdateLayoutConfig
	BackgroundConfig           []api.BackgroundConfig
}
