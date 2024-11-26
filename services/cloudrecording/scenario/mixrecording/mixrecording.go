package mixrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireMixRecodingClientRequest struct {
	ResourceExpiredHour int
	ExcludeResourceIds  []string
	RegionAffinity      int

	// StartParameter 设置该字段后，可以提升可用性并优化负载均衡。
	//
	// 注意：如果填写该字段，则必须确保 startParameter object 和后续 start 请求中填写的 clientRequest object 完全一致，
	// 且取值合法，否则 start 请求会收到报错。
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
