package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Acquire struct {
	baseAPI *api.Acquire
}

func NewAcquire(baseAPI *api.Acquire) *Acquire {
	return &Acquire{baseAPI: baseAPI}
}

func (a *Acquire) Do(ctx context.Context, cname string, uid string, enablePostpone bool, clientRequest *AcquireIndividualRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			SnapshotConfig:      clientRequest.StartParameter.SnapshotConfig,
			AppsCollection:      clientRequest.StartParameter.AppsCollection,
			TranscodeOptions:    clientRequest.StartParameter.TranscodeOptions,
		}
	}

	scene := 0
	if enablePostpone {
		scene = 2
	}
	return a.baseAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               scene,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}
