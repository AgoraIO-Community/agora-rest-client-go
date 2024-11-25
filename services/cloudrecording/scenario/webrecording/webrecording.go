package webrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireWebRecodingClientRequest struct {
	ResourceExpiredHour int      `json:"resourceExpiredHour,omitempty"`
	ExcludeResourceIds  []string `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int      `json:"regionAffinity,omitempty"`

	// StartParameter 设置该字段后，可以提升可用性并优化负载均衡。
	//
	// 注意：如果填写该字段，则必须确保 startParameter object 和后续 start 请求中填写的 clientRequest object 完全一致，
	// 且取值合法，否则 start 请求会收到报错。
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
