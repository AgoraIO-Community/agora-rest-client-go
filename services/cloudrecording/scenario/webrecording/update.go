package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Update struct {
	baseAPI *api.Update
}

func NewUpdate(baseUpdate *api.Update) *Update {
	return &Update{baseAPI: baseUpdate}
}

func (s *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateWebRecordingClientRequest) (*api.UpdateResp, error) {
	return s.baseAPI.Do(ctx, resourceID, sid, api.WebMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			WebRecordingConfig: clientRequest.WebRecordingConfig,
			RtmpPublishConfig:  clientRequest.RtmpPublishConfig,
		},
	})
}
