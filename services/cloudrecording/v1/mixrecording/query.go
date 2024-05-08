package mixrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Query struct {
	Base *baseV1.Query
}

var _ baseV1.QueryMixRecording = (*Query)(nil)

func (q Query) DoHLS(ctx context.Context, resourceID string, sid string) (*baseV1.QueryMixRecordingHLSResp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, baseV1.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.QueryMixRecordingHLSResp

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.QueryMixRecordingHLSSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

func (q Query) DoHLSAndMP4(ctx context.Context, resourceID string, sid string) (*baseV1.QueryMixRecordingHLSAndMP4Resp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, baseV1.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp baseV1.QueryMixRecordingHLSAndMP4Resp

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		mixResp.SuccessResp = baseV1.QueryMixRecordingHLSAndMP4SuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
