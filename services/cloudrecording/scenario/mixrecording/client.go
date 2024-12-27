package mixrecording

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

func (c *Client) Acquire(ctx context.Context, cname string, uid string,
	clientRequest *AcquireMixRecodingClientRequest) (*api.AcquireResp, error) {

	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
		}
	}

	return c.acquireAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               0,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}

func (c *Client) QueryHLS(ctx context.Context, resourceID string, sid string,
) (*QueryMixRecordingHLSResp, error) {

	resp, err := c.queryAPI.Do(ctx, resourceID, sid, api.MixMode)
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

func (c *Client) QueryHLSAndMP4(ctx context.Context, resourceID string, sid string,
) (*QueryMixRecordingHLSAndMP4Resp, error) {

	resp, err := c.queryAPI.Do(ctx, resourceID, sid, api.MixMode)
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
func (c *Client) Start(ctx context.Context, resourceID string, cname string, uid string,
	clientRequest *StartMixRecordingClientRequest) (*api.StartResp, error) {

	return c.startAPI.Do(ctx, resourceID, api.MixMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			RecordingConfig:     clientRequest.RecordingConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}

func (c *Client) Stop(ctx context.Context, resourceID string, sid string, cname string, uid string,
	asyncStop bool) (*api.StopResp, error) {

	return c.stopAPI.Do(ctx, resourceID, sid, api.MixMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}

func (c *Client) Update(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *UpdateMixRecordingClientRequest) (*api.UpdateResp, error) {

	return c.updateAPI.Do(ctx, resourceID, sid, api.MixMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}

func (c *Client) UpdateLayout(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *UpdateLayoutUpdateMixRecordingClientRequest) (*api.UpdateLayoutResp, error) {

	return c.updateLayoutAPI.Do(ctx, resourceID, sid, api.MixMode, &api.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateLayoutClientRequest{
			MaxResolutionUID:           clientRequest.MaxResolutionUID,
			MixedVideoLayout:           clientRequest.MixedVideoLayout,
			BackgroundColor:            clientRequest.BackgroundColor,
			BackgroundImage:            clientRequest.BackgroundImage,
			DefaultUserBackgroundImage: clientRequest.DefaultUserBackgroundImage,
			LayoutConfig:               clientRequest.LayoutConfig,
			BackgroundConfig:           clientRequest.BackgroundConfig,
		},
	})
}
