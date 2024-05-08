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

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, payload *baseV1.StopReqBody) (*baseV1.StopWebRecordingResp, error) {
	mode := baseV1.WebMode
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, mode, payload)
	if err != nil {
		return nil, err
	}

	var webResp baseV1.StopWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		webResp.SuccessResp = baseV1.StopWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetWebRecordingServerResponse(),
		}
	}

	return &webResp, nil
}
