package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Start struct {
	baseAPI *api.Start
}

func NewStart(base *api.Start) *Start {
	return &Start{baseAPI: base}
}

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartWebRecordingClientRequest) (*api.StartResp, error) {
	return s.baseAPI.Do(ctx, resourceID, api.WebMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.RecordingFileConfig,
			StorageConfig:          clientRequest.StorageConfig,
			ExtensionServiceConfig: clientRequest.ExtensionServiceConfig,
		},
	})
}
