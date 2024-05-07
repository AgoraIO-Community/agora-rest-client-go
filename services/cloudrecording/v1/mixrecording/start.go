package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Start struct {
	Base *baseV1.Start
}

var _ baseV1.StartMixRecording = (*Start)(nil)

func (s *Start) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.StartMixRecording {
	s.Base.WithForwardRegion(prefix)

	return s
}

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartMixRecordingClientRequest) (*baseV1.StartResp, error) {
	return s.Base.Do(ctx, resourceID, baseV1.MixMode, &baseV1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StartClientRequest{
			Token:               clientRequest.Token,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			RecordingConfig:     clientRequest.RecordingConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}
