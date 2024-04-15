package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Starter struct {
	BaseStarter *baseV1.Starter
}

var _ baseV1.StartMixRecording = (*Starter)(nil)

func (s *Starter) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartMixRecordingClientRequest) (*baseV1.StarterResp, error) {
	return s.BaseStarter.Do(ctx, resourceID, baseV1.MixMode, &baseV1.StartReqBody{
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
