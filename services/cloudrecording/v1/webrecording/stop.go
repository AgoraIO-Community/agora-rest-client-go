package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Stop struct {
	BaseStop *baseV1.Stop
}

var _ baseV1.StopWebRecording = (*Stop)(nil)

func (s *Stop) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.StopWebRecording {
	s.BaseStop.WithForwardRegion(prefix)

	return s
}

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*baseV1.StopWebRecordingResp, error) {
	mode := baseV1.WebMode
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, mode, &baseV1.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var webResp baseV1.StopWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		webResp.SuccessResponse = baseV1.StopWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetWebRecordingServerResponse(),
		}
	}

	return &webResp, nil
}
