package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Query struct {
	baseAPI *api.Query
}

func NewQuery(base *api.Query) *Query {
	return &Query{baseAPI: base}
}

func (q *Query) Do(ctx context.Context, resourceID string, sid string) (*QueryWebRecordingResp, error) {
	resp, err := q.baseAPI.Do(ctx, resourceID, sid, api.WebMode)
	if err != nil {
		return nil, err
	}

	var webResp QueryWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		webResp.SuccessResponse = QueryWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetWebRecording2CDNServerResponse(),
		}
	}

	return &webResp, nil
}
