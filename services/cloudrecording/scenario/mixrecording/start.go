package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Start struct {
	Base *api.Start
}

func NewStart(base *api.Start) *Start {
	return &Start{Base: base}
}

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartMixRecordingClientRequest) (*api.StartResp, error) {
	return s.Base.Do(ctx, resourceID, api.MixMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			RecordingConfig:     clientRequest.RecordingConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}
