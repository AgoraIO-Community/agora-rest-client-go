package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Start struct {
	Base *baseV1.Start
}

var _ baseV1.StartWebRecording = (*Start)(nil)

func (s *Start) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.StartWebRecording {
	s.Base.WithForwardRegion(prefix)

	return s
}

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartWebRecordingClientRequest) (*baseV1.StarterResp, error) {
	return s.Base.Do(ctx, resourceID, baseV1.WebMode, &baseV1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StartClientRequest{
			RecordingFileConfig:    clientRequest.RecordingFileConfig,
			StorageConfig:          clientRequest.StorageConfig,
			ExtensionServiceConfig: clientRequest.ExtensionServiceConfig,
		},
	})
}
