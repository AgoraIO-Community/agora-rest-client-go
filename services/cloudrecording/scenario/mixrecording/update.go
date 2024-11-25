package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Update struct {
	Base *api.Update
}

func NewUpdate(base *api.Update) *Update {
	return &Update{Base: base}
}

func (u *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *UpdateMixRecordingClientRequest,
) (*api.UpdateResp, error) {
	return u.Base.Do(ctx, resourceID, sid, api.MixMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}
