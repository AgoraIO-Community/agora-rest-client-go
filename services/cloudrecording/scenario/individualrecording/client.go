package individualrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type Client struct {
	acquireAPI      *api.Acquire
	startAPI        *api.Start
	stopAPI         *api.Stop
	queryAPI        *api.Query
	updateLayoutAPI *api.UpdateLayout
	updateAPI       *api.Update
}

func NewClient(
	acquireAPI *api.Acquire,
	startAPI *api.Start,
	stopAPI *api.Stop,
	queryAPI *api.Query,
	updateLayoutAPI *api.UpdateLayout,
	updateAPI *api.Update,
) *Client {
	return &Client{
		acquireAPI:      acquireAPI,
		startAPI:        startAPI,
		stopAPI:         stopAPI,
		queryAPI:        queryAPI,
		updateLayoutAPI: updateLayoutAPI,
		updateAPI:       updateAPI,
	}
}

func (c *Client) Acquire(ctx context.Context, cname string, uid string, enablePostpone bool, clientRequest *AcquireIndividualRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			SnapshotConfig:      clientRequest.StartParameter.SnapshotConfig,
			AppsCollection:      clientRequest.StartParameter.AppsCollection,
			TranscodeOptions:    clientRequest.StartParameter.TranscodeOptions,
		}
	}

	scene := 0
	if enablePostpone {
		scene = 2
	}
	return c.acquireAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               scene,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}

func (c *Client) Start(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartIndividualRecordingClientRequest) (*api.StartResp, error) {
	return c.startAPI.Do(ctx, resourceID, api.IndividualMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			AppsCollection:      clientRequest.AppsCollection,
			RecordingConfig:     clientRequest.RecordingConfig,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			TranscodeOptions:    clientRequest.TranscodeOptions,
			SnapshotConfig:      clientRequest.SnapshotConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}

func (c *Client) Query(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingResp, error) {
	resp, err := c.queryAPI.Do(ctx, resourceID, sid, api.IndividualMode)
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

func (c *Client) QueryVideoScreenshot(ctx context.Context, resourceID string, sid string) (*QueryIndividualRecordingVideoScreenshotResp, error) {
	resp, err := c.queryAPI.Do(ctx, resourceID, sid, api.IndividualMode)
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

func (c *Client) Update(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateIndividualRecordingClientRequest) (*api.UpdateResp, error) {
	return c.updateAPI.Do(ctx, resourceID, sid, api.IndividualMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}

func (c *Client) Stop(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*api.StopResp, error) {
	return c.stopAPI.Do(ctx, resourceID, sid, api.IndividualMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}
