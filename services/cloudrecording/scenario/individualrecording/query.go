package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Query struct {
	baseAPI *api.Query
}

func NewQuery(baseAPI *api.Query) *Query {
	return &Query{baseAPI: baseAPI}
}

func (q *Query) Do(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingResp, error) {
	resp, err := q.baseAPI.Do(ctx, resourceID, sid, api.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp QueryIndividualRecordingResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = QueryIndividualRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualRecordingServerResponse(),
		}
	}

	return &individualResp, nil
}

func (q *Query) DoVideoScreenshot(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingVideoScreenshotResp, error) {
	resp, err := q.baseAPI.Do(ctx, resourceID, sid, api.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp QueryIndividualRecordingVideoScreenshotResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResponse
		individualResp.SuccessResponse = QueryIndividualRecordingVideoScreenshotSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualVideoScreenshotServerResponse(),
		}
	}

	return &individualResp, nil
}
