package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type UpdateLayout struct {
	BaseUpdateLayout *baseV1.UpdateLayout
}

var _ baseV1.UpdateLayoutMixRecording = (*UpdateLayout)(nil)

func (u *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *baseV1.UpdateLayoutUpdateMixRecordingClientRequest) (*baseV1.UpdateLayoutResp, error) {
	return u.BaseUpdateLayout.Do(ctx, resourceID, sid, baseV1.MixMode, &baseV1.UpdateLayoutReqBody{
		Cname:         cname,
		Uid:           uid,
		ClientRequest: &baseV1.UpdateLayoutClientRequest{},
	})
}
