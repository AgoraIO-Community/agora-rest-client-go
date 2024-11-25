package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Update struct {
	baseAPI *api.Update
}

func NewUpdate(base *api.Update) *Update {
	return &Update{baseAPI: base}
}

func (u *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateIndividualRecordingClientRequest) (*api.UpdateResp, error) {
	return u.baseAPI.Do(ctx, resourceID, sid, api.IndividualMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}
