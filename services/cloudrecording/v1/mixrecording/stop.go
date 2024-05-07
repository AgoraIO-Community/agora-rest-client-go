package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Stop struct {
	Base *baseV1.Stop
}

var _ baseV1.StopMixRecording = (*Stop)(nil)

func (s *Stop) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.StopMixRecording {
	s.Base.WithForwardRegion(prefix)

	return s
}

func (s *Stop) DoHLS(ctx context.Context, resourceID string, sid string, payload *baseV1.StopReqBody) (*baseV1.StopMixRecordingHLSSuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, baseV1.MixMode, payload)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSSuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.StopMixRecordingHLSResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

func (s *Stop) DoHLSAndMP4(ctx context.Context, resourceID string, sid string, payload *baseV1.StopReqBody) (*baseV1.StopMixRecordingHLSAndMP4SuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, baseV1.MixMode, payload)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSAndMP4SuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.StopMixRecordingHLSAndMP4Resp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
