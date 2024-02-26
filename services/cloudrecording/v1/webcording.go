package v1

import (
	"context"
)

type AcquirerWebRecodingClientRequest struct {
	ResourceExpiredHour int      `json:"resourceExpiredHour"`
	ExcludeResourceIds  []string `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int      `json:"regionAffinity,omitempty"`
}

type AcquireWebRecording interface {
	Do(ctx context.Context, cname string, uid string, clientRequest *AcquirerWebRecodingClientRequest) (*AcquirerResp, error)
}

type QueryWebRecordingResp struct {
	Response
	SuccessResp QueryWebRecordingSuccessResp
}

type QueryWebRecordingSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse QueryWebRecordingServerResponse
}
type QueryWebRecording interface {
	Do(ctx context.Context, resourceID string, sid string) (*QueryWebRecordingResp, error)
}

type StartWebRecordingClientRequest struct {
	AppsCollection         *AppsCollection         `json:"appsCollection,omitempty"`
	RecordingFileConfig    *RecordingFileConfig    `json:"recordingFileConfig,omitempty"`
	StorageConfig          *StorageConfig          `json:"storageConfig,omitempty"`
	ExtensionServiceConfig *ExtensionServiceConfig `json:"extensionServiceConfig,omitempty"`
}

type StartWebRecording interface {
	Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartWebRecordingClientRequest) (*StarterResp, error)
}

type StopWebRecordingResp struct {
	Response
	SuccessResp StopWebRecordingSuccessResp
}

type StopWebRecordingSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse StopWebRecordingServerResponse
}

type StopWebRecording interface {
	Do(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopWebRecordingResp, error)
}

type UpdateWebRecordingClientRequest struct {
	WebRecordingConfig *UpdateWebRecordingConfig `json:"webRecordingConfig,omitempty"`
	RtmpPublishConfig  *UpdateRtmpPublishConfig  `json:"rtmpPublishConfig,omitempty"`
}

type UpdateWebRecording interface {
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
