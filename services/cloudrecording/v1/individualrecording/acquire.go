package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Acquire struct {
	Base *baseV1.Acquire
}

var _ baseV1.AcquireIndividualRecording = (*Acquire)(nil)

func (a *Acquire) Do(ctx context.Context, cname string, uid string, enablePostponeTranscodingMix bool, clientRequest *baseV1.AcquirerIndividualRecodingClientRequest) (*baseV1.AcquirerResp, error) {
	scene := 0
	if enablePostponeTranscodingMix {
		scene = 2
	}
	return a.Base.Do(ctx, &baseV1.AcquirerReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.AcquirerClientRequest{
			Scene:               scene,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
		},
	})
}
