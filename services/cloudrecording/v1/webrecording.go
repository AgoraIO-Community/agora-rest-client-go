package v1

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
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

type AcquireWebRecording interface {
	WithForwardRegion(prefix core.ForwardedReginPrefix) AcquireWebRecording
	Do(ctx context.Context, cname string, uid string, clientRequest *AcquireWebRecodingClientRequest) (*AcquireResp, error)
}

type QueryWebRecordingResp struct {
	Response
	SuccessResp QueryWebRecordingSuccessResp
}

type QueryWebRecordingSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *QueryWebRecordingServerResponse
}
type QueryWebRecording interface {
	WithForwardRegion(prefix core.ForwardedReginPrefix) QueryWebRecording
	Do(ctx context.Context, resourceID string, sid string) (*QueryWebRecordingResp, error)
}

type StartWebRecordingClientRequest struct {
	RecordingFileConfig    *RecordingFileConfig
	StorageConfig          *StorageConfig
	ExtensionServiceConfig *ExtensionServiceConfig
}

type StartWebRecording interface {
	WithForwardRegion(prefix core.ForwardedReginPrefix) StartWebRecording
	Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartWebRecordingClientRequest) (*StartResp, error)
}

type StopWebRecordingResp struct {
	Response
	SuccessResp StopWebRecordingSuccessResp
}

type StopWebRecordingSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *StopWebRecordingServerResponse
}

type StopWebRecording interface {
	WithForwardRegion(prefix core.ForwardedReginPrefix) StopWebRecording
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopWebRecordingResp, error)
}

type UpdateWebRecordingClientRequest struct {
	WebRecordingConfig *UpdateWebRecordingConfig
	RtmpPublishConfig  *UpdateRtmpPublishConfig
}

type UpdateWebRecording interface {
	WithForwardRegion(prefix core.ForwardedReginPrefix) UpdateWebRecording
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateWebRecordingClientRequest) (*UpdateResp, error)
}

type WebRecording interface {
	SetBase(base *BaseCollection)
	Acquire() AcquireWebRecording
	Query() QueryWebRecording
	Start() StartWebRecording
	Stop() StopWebRecording
	Update() UpdateWebRecording
}
