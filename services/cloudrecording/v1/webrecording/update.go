package webrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Update struct {
	BaseUpdate *baseV1.Update
}

var _ baseV1.UpdateWebRecording = (*Update)(nil)

func (s *Update) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *baseV1.UpdateWebRecordingClientRequest) (*baseV1.UpdateResp, error) {
	return s.BaseUpdate.Do(ctx, resourceID, sid, baseV1.WebMode, &baseV1.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.UpdateClientRequest{
			WebRecordingConfig: clientRequest.WebRecordingConfig,
			RtmpPublishConfig:  clientRequest.RtmpPublishConfig,
		},
	})
}
