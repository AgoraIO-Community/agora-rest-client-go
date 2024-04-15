package v1

import (
	"context"
)

type AcquirerMixRecodingClientRequest struct {
	ResourceExpiredHour int      `json:"resourceExpiredHour"`
	ExcludeResourceIds  []string `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int      `json:"regionAffinity,omitempty"`
}

type AcquireMixRecording interface {
	Do(ctx context.Context, cname string, uid string, clientRequest *AcquirerMixRecodingClientRequest) (*AcquirerResp, error)
}

type QueryMixRecordingHLSSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse QueryMixRecordingHLSServerResponse
}

type QueryMixRecordingHLSResp struct {
	Response
	SuccessResp QueryMixRecordingHLSSuccessResp
}

type QueryMixRecordingHLSAndMP4SuccessResp struct {
	ResourceId     string
	SID            string
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
	SID            string
	ServerResponse StopMixRecordingHLSServerResponse
}

type StopMixRecordingHLSSuccessResponse struct {
	Response
	SuccessResp StopMixRecordingHLSResp
}

type StopMixRecordingHLSAndMP4Resp struct {
	ResourceId     string
	SID            string
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
}
type UpdateMixRecording interface {
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateMixRecordingClientRequest) (*UpdateResp, error)
}

type UpdateLayoutUpdateMixRecordingClientRequest struct {
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
