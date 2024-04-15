package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Update struct {
	BaseUpdate *baseV1.Update
}

var _ baseV1.UpdateMixRecording = (*Update)(nil)

func (i *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *baseV1.UpdateMixRecordingClientRequest) (*baseV1.UpdateResp, error) {
	return i.BaseUpdate.Do(ctx, resourceID, sid, baseV1.MixMode, &baseV1.UpdateReqBody{
		Cname:         cname,
		Uid:           uid,
		ClientRequest: &baseV1.UpdateClientRequest{},
	})
}
