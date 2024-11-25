package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Stop struct {
	baseAPI *api.Stop
}

func NewStop(baseAPI *api.Stop) *Stop {
	return &Stop{baseAPI: baseAPI}
}

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopIndividualRecordingResp, error) {
	resp, err := s.baseAPI.Do(ctx, resourceID, sid, api.IndividualMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var individualResp StopIndividualRecordingResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = StopIndividualRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualRecordingServerResponse(),
		}
	}

	return &individualResp, nil
}

func (s *Stop) DoVideoScreenshot(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopIndividualRecordingVideoScreenshotResp, error) {
	resp, err := s.baseAPI.Do(ctx, resourceID, sid, api.IndividualMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var individualResp StopIndividualRecordingVideoScreenshotResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = StopIndividualRecordingVideoScreenshotSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualVideoScreenshotServerResponse(),
		}
	}

	return &individualResp, nil
}
