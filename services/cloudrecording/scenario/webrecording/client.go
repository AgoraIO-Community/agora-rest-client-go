package webrecording

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

func (c *Client) Acquire(ctx context.Context, cname string, uid string, clientRequest *AcquireWebRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:          clientRequest.StartParameter.StorageConfig,
			ExtensionServiceConfig: clientRequest.StartParameter.ExtensionServiceConfig,
		}
	}

	return c.acquireAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               1,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}

func (c *Client) Query(ctx context.Context, resourceID string, sid string) (*QueryWebRecordingResp, error) {
	resp, err := c.queryAPI.Do(ctx, resourceID, sid, api.WebMode)
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

func (c *Client) Start(ctx context.Context, resourceID string, cname string, uid string, clientRequest *StartWebRecordingClientRequest) (*api.StartResp, error) {
	return c.startAPI.Do(ctx, resourceID, api.WebMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.RecordingFileConfig,
			StorageConfig:          clientRequest.StorageConfig,
			ExtensionServiceConfig: clientRequest.ExtensionServiceConfig,
		},
	})
}

func (c *Client) Stop(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*api.StopResp, error) {
	return c.stopAPI.Do(ctx, resourceID, sid, api.WebMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}

func (c *Client) Update(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateWebRecordingClientRequest) (*api.UpdateResp, error) {
	return c.updateAPI.Do(ctx, resourceID, sid, api.WebMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			WebRecordingConfig: clientRequest.WebRecordingConfig,
			RtmpPublishConfig:  clientRequest.RtmpPublishConfig,
		},
	})
}
