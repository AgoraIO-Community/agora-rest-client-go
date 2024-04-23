package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Starter struct {
	Base *baseV1.Starter
}

var _ baseV1.StartIndividualRecording = (*Starter)(nil)

func (s *Starter) Do(ctx context.Context, resourceID string, cname string, uid string, clientRequest *baseV1.StartIndividualRecordingClientRequest) (*baseV1.StarterResp, error) {
	return s.Base.Do(ctx, resourceID, baseV1.IndividualMode, &baseV1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StartClientRequest{
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
