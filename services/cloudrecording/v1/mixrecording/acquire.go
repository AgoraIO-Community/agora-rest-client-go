package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Acquire struct {
	Base *baseV1.Acquire
}

var _ baseV1.AcquireMixRecording = (*Acquire)(nil)

func (a *Acquire) Do(ctx context.Context, cname string, uid string, clientRequest *baseV1.AcquirerMixRecodingClientRequest) (*baseV1.AcquirerResp, error) {
	var startParameter *baseV1.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &baseV1.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
		}
	}

	return a.Base.Do(ctx, &baseV1.AcquirerReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.AcquirerClientRequest{
			Scene:               0,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}