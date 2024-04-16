package v1

import (
	"context"
)

type AcquirerIndividualRecodingClientRequest struct {
	ResourceExpiredHour int      `json:"resourceExpiredHour"`
	ExcludeResourceIds  []string `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int      `json:"regionAffinity,omitempty"`
}

type AcquireIndividualRecording interface {
	Do(ctx context.Context, cname string, uid string, enablePostponeTranscodingMix bool, clientRequest *AcquirerIndividualRecodingClientRequest) (*AcquirerResp, error)
}

type StartIndividualRecordingClientRequest struct {
	Token               string
	StorageConfig       *StorageConfig
	RecordingConfig     *RecordingConfig
	RecordingFileConfig *RecordingFileConfig
	SnapshotConfig      *SnapshotConfig
	AppsCollection      *AppsCollection
	TranscodeOptions    *TranscodeOptions
}

type StartIndividualRecording interface {
	Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartIndividualRecordingClientRequest) (*StarterResp, error)
}

type QueryIndividualRecordingSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse QueryIndividualRecordingServerResponse
}

type QueryIndividualRecordingResp struct {
	Response
	SuccessResp QueryIndividualRecordingSuccessResp
}

type QueryIndividualRecordingVideoScreenshotSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse QueryIndividualVideoScreenshotServerResponse
}

type QueryIndividualRecordingVideoScreenshotResp struct {
	Response
	SuccessResp QueryIndividualRecordingVideoScreenshotSuccessResp
}

type QueryIndividualRecording interface {
	Do(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingResp, error)
	DoVideoScreenshot(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingVideoScreenshotResp, error)
}

type UpdateIndividualRecordingClientRequest struct {
	StreamSubscribe *UpdateStreamSubscribe
}

type UpdateIndividualRecording interface {
	Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateIndividualRecordingClientRequest) (*UpdateResp, error)
}

type StopIndividualRecordingSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse StopIndividualRecordingServerResponse
}

type StopIndividualRecordingVideoScreenshotSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse StopIndividualVideoScreenshotServerResponse
}

type StopIndividualRecordingResp struct {
	Response
	SuccessResp StopIndividualRecordingSuccessResp
}

type StopIndividualRecordingVideoScreenshotResp struct {
	Response
	SuccessResp StopIndividualRecordingVideoScreenshotSuccessResp
}

type StopIndividualRecording interface {
	Do(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopIndividualRecordingResp, error)
	DoVideoScreenshot(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopIndividualRecordingVideoScreenshotResp, error)
}

type IndividualRecording interface {
	SetBase(base *BaseCollection)
	Acquire() AcquireIndividualRecording
	Start() StartIndividualRecording
	Query() QueryIndividualRecording
	Update() UpdateIndividualRecording
	Stop() StopIndividualRecording
}
