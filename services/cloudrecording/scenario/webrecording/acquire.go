package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Acquire struct {
	baseAPI *api.Acquire
}

func NewAcquire(base *api.Acquire) *Acquire {
	return &Acquire{baseAPI: base}
}

func (a *Acquire) Do(ctx context.Context, cname string, uid string, clientRequest *AcquireWebRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:          clientRequest.StartParameter.StorageConfig,
			ExtensionServiceConfig: clientRequest.StartParameter.ExtensionServiceConfig,
		}
	}

	return a.baseAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               1,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}
