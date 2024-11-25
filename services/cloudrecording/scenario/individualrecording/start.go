package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Start struct {
	baseAPI *api.Start
}

func NewStart(baseAPI *api.Start) *Start {
	return &Start{baseAPI: baseAPI}
}

func (s *Start) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartIndividualRecordingClientRequest) (*api.StartResp, error) {
	return s.baseAPI.Do(ctx, resourceID, api.IndividualMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			AppsCollection:      clientRequest.AppsCollection,
			RecordingConfig:     clientRequest.RecordingConfig,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			TranscodeOptions:    clientRequest.TranscodeOptions,
			SnapshotConfig:      clientRequest.SnapshotConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}
