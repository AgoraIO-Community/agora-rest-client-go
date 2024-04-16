package webrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Starter struct {
	Base *baseV1.Starter
}

var _ baseV1.StartWebRecording = (*Starter)(nil)

func (s *Starter) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartWebRecordingClientRequest) (*baseV1.StarterResp, error) {
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
