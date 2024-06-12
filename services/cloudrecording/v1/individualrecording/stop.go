package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Stop struct {
	BaseStop *baseV1.Stop
}

var _ baseV1.StopIndividualRecording = (*Stop)(nil)

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*baseV1.StopIndividualRecordingResp, error) {
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, baseV1.IndividualMode, &baseV1.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var individualResp baseV1.StopIndividualRecordingResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = baseV1.StopIndividualRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualRecordingServerResponse(),
		}
	}

	return &individualResp, nil
}

func (s *Stop) DoVideoScreenshot(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*baseV1.StopIndividualRecordingVideoScreenshotResp, error) {
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, baseV1.IndividualMode, &baseV1.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var individualResp baseV1.StopIndividualRecordingVideoScreenshotResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = baseV1.StopIndividualRecordingVideoScreenshotSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualVideoScreenshotServerResponse(),
		}
	}

	return &individualResp, nil
}
