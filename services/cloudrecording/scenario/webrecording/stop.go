package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Stop struct {
	baseAPI *api.Stop
}

func NewStop(baseStop *api.Stop) *Stop {
	return &Stop{baseAPI: baseStop}
}

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopWebRecordingResp, error) {
	mode := api.WebMode
	resp, err := s.baseAPI.Do(ctx, resourceID, sid, mode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var webResp StopWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		webResp.SuccessResponse = StopWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetWebRecordingServerResponse(),
		}
	}

	return &webResp, nil
}
