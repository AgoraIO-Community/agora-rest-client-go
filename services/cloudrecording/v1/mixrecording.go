package v1

import (
	"context"
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

type AcquireMixRecording interface {
	Do(ctx context.Context, cname string, uid string, clientRequest *AcquireMixRecodingClientRequest) (*AcquireResp, error)
}

type QueryMixRecordingHLSSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse QueryMixRecordingHLSServerResponse
}

type QueryMixRecordingHLSResp struct {
	Response
	SuccessResp QueryMixRecordingHLSSuccessResp
}

type QueryMixRecordingHLSAndMP4SuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse QueryMixRecordingHLSAndMP4ServerResponse
}
type QueryMixRecordingHLSAndMP4Resp struct {
	Response
	SuccessResp QueryMixRecordingHLSAndMP4SuccessResp
}

type QueryMixRecording interface {
	DoHLS(ctx context.Context, resourceID string, sid string) (*QueryMixRecordingHLSResp, error)
	DoHLSAndMP4(ctx context.Context, resourceID string, sid string) (*QueryMixRecordingHLSAndMP4Resp, error)
}

type StartMixRecordingClientRequest struct {
	Token               string
	RecordingConfig     *RecordingConfig
	RecordingFileConfig *RecordingFileConfig
	StorageConfig       *StorageConfig
}

type StartMixRecording interface {
	Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartMixRecordingClientRequest) (*StarterResp, error)
}

type StopMixRecordingHLSResp struct {
	ResourceId     string
	Sid            string
	ServerResponse StopMixRecordingHLSServerResponse
}

type StopMixRecordingHLSSuccessResponse struct {
	Response
	SuccessResp StopMixRecordingHLSResp
}

type StopMixRecordingHLSAndMP4Resp struct {
	ResourceId     string
	Sid            string
	ServerResponse StopMixRecordingHLSAndMP4ServerResponse
}

type StopMixRecordingHLSAndMP4SuccessResponse struct {
	Response
	SuccessResp StopMixRecordingHLSAndMP4Resp
}

type StopMixRecording interface {
	DoHLS(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopMixRecordingHLSSuccessResponse, error)
	DoHLSAndMP4(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopMixRecordingHLSAndMP4SuccessResponse, error)
}

type UpdateMixRecordingClientRequest struct {
	StreamSubscribe *UpdateStreamSubscribe
}
type UpdateMixRecording interface {
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateMixRecordingClientRequest) (*UpdateResp, error)
}

type UpdateLayoutUpdateMixRecordingClientRequest struct {
	MaxResolutionUID           string
	MixedVideoLayout           int
	BackgroundColor            string
	BackgroundImage            string
	DefaultUserBackgroundImage string
	LayoutConfig               []UpdateLayoutConfig
	BackgroundConfig           []BackgroundConfig
}

type UpdateLayoutMixRecording interface {
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateLayoutUpdateMixRecordingClientRequest) (*UpdateLayoutResp, error)
}

type MixRecording interface {
	SetBase(base *BaseCollection)
	Acquire() AcquireMixRecording
	Query() QueryMixRecording
	Start() StartMixRecording
	Stop() StopMixRecording
	Update() UpdateMixRecording
	UpdateLayout() UpdateLayoutMixRecording
}
