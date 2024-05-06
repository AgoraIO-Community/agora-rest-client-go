package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Query struct {
	Base *baseV1.Query
}

var _ baseV1.QueryIndividualRecording = (*Query)(nil)

func (q *Query) Do(ctx context.Context, resourceID string, sid string) (*baseV1.QueryIndividualRecordingResp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, baseV1.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp baseV1.QueryIndividualRecordingResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		individualResp.SuccessResp = baseV1.QueryIndividualRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetIndividualRecordingServerResponse(),
		}
	}

	return &individualResp, nil
}

func (q *Query) DoVideoScreenshot(ctx context.Context, resourceID string, sid string) (*baseV1.QueryIndividualRecordingVideoScreenshotResp, error) {
	resp, err := q.Base.Do(ctx, resourceID, sid, baseV1.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp baseV1.QueryIndividualRecordingVideoScreenshotResp

	individualResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		individualResp.SuccessResp = baseV1.QueryIndividualRecordingVideoScreenshotSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: *successResp.GetIndividualVideoScreenshotServerResponse(),
		}
	}

	return &individualResp, nil
}
