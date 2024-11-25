package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Acquire struct {
	base *api.Acquire
}

func NewAcquire(base *api.Acquire) *Acquire {
	return &Acquire{base: base}
}

func (a *Acquire) Do(ctx context.Context, cname string, uid string, clientRequest *AcquireMixRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
		}
	}

	return a.base.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               0,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}
