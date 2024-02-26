package webrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Acquire struct {
	BaseAcquire *baseV1.Acquire
}

var _ baseV1.AcquireWebRecording = (*Acquire)(nil)

func (a *Acquire) Do(ctx context.Context, cname string, uid string, clientRequest *baseV1.AcquirerWebRecodingClientRequest) (*baseV1.AcquirerResp, error) {
	return a.BaseAcquire.Do(ctx, &baseV1.AcquirerReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.AcquirerClientRequest{
			Scene:               1,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
		},
	})
}
