package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Update struct {
	Base *baseV1.Update
}

var _ baseV1.UpdateIndividualRecording = (*Update)(nil)

func (u *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *baseV1.UpdateIndividualRecordingClientRequest) (*baseV1.UpdateResp, error) {
	return u.Base.Do(ctx, resourceID, sid, baseV1.IndividualMode, &baseV1.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}
