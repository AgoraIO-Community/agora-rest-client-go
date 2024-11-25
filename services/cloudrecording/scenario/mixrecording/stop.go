package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Stop struct {
	Base *api.Stop
}

func NewStop(base *api.Stop) *Stop {
	return &Stop{Base: base}
}

func (s *Stop) DoHLS(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopMixRecordingHLSSuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, api.MixMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var mixResp StopMixRecordingHLSSuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = StopMixRecordingHLSResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

func (s *Stop) DoHLSAndMP4(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*StopMixRecordingHLSAndMP4SuccessResponse, error) {
	resp, err := s.Base.Do(ctx, resourceID, sid, api.MixMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
	if err != nil {
		return nil, err
	}

	var mixResp StopMixRecordingHLSAndMP4SuccessResponse

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = StopMixRecordingHLSAndMP4Resp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
