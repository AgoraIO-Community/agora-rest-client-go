package webrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Start struct {
	Base *baseV1.Start
}

var _ baseV1.StartWebRecording = (*Start)(nil)

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartWebRecordingClientRequest) (*baseV1.StartResp, error) {
	mode := baseV1.WebMode
	return s.Base.Do(ctx, resourceID, mode, &baseV1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StartClientRequest{
			RecordingFileConfig:    clientRequest.RecordingFileConfig,
			StorageConfig:          clientRequest.StorageConfig,
			ExtensionServiceConfig: clientRequest.ExtensionServiceConfig,
		},
	})
}
