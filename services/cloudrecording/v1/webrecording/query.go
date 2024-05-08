package webrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Query struct {
	Base *baseV1.Query
}

var _ baseV1.QueryWebRecording = (*Query)(nil)

func (q *Query) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.QueryWebRecording {
	q.Base.WithForwardRegion(prefix)

	return q
}

func (q *Query) Do(ctx context.Context, resourceID string, sid string) (*baseV1.QueryWebRecordingResp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, baseV1.WebMode)
	if err != nil {
		return nil, err
	}

	var webResp baseV1.QueryWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		webResp.SuccessResp = baseV1.QueryWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetWebRecording2CDNServerResponse(),
		}
	}

	return &webResp, nil
}
