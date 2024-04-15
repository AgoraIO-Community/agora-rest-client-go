package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Stop struct {
	BaseStop *baseV1.Stop
}

var _ baseV1.StopMixRecording = (*Stop)(nil)

func (s *Stop) DoHLS(ctx context.Context, resourceID string, sid string, payload *baseV1.StopReqBody) (*baseV1.StopMixRecordingHLSSuccessResponse, error) {
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, baseV1.MixMode, payload)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSSuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.StopMixRecordingHLSResp{
			ResourceId:     successResp.ResourceId,
			SID:            successResp.SID,
			ServerResponse: *successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil

}

func (s *Stop) DoHLSAndMP4(ctx context.Context, resourceID string, sid string, payload *baseV1.StopReqBody) (*baseV1.StopMixRecordingHLSAndMP4SuccessResponse, error) {
	resp, err := s.BaseStop.Do(ctx, resourceID, sid, baseV1.MixMode, payload)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSAndMP4SuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.StopMixRecordingHLSAndMP4Resp{
			ResourceId:     successResp.ResourceId,
			SID:            successResp.SID,
			ServerResponse: *successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
