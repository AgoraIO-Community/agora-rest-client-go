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

func (s *Stop) DoHLS(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*baseV1.StopMixRecordingHLSSuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, baseV1.MixMode, &baseV1.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSSuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = baseV1.StopMixRecordingHLSResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

func (s *Stop) DoHLSAndMP4(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*baseV1.StopMixRecordingHLSAndMP4SuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, baseV1.MixMode, &baseV1.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.StopMixRecordingHLSAndMP4SuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = baseV1.StopMixRecordingHLSAndMP4Resp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
