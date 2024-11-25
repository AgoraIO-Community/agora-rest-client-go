package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Query struct {
	Base *api.Query
}

func NewQuery(base *api.Query) *Query {
	return &Query{Base: base}
}

func (q *Query) DoHLS(ctx context.Context, resourceID string, sid string) (*QueryMixRecordingHLSResp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, api.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp QueryMixRecordingHLSResp

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = QueryMixRecordingHLSSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

func (q *Query) DoHLSAndMP4(ctx context.Context, resourceID string, sid string) (*QueryMixRecordingHLSAndMP4Resp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, api.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp QueryMixRecordingHLSAndMP4Resp

	mixResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		mixResp.SuccessResponse = QueryMixRecordingHLSAndMP4SuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}
