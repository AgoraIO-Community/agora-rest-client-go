package cloudrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	agoraClient "github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/individualrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/mixrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/webrecording"
)

const projectName = "cloud_recording"

type Client struct {
	acquireAPI      *api.Acquire
	startAPI        *api.Start
	stopAPI         *api.Stop
	queryAPI        *api.Query
	updateLayoutAPI *api.UpdateLayout
	updateAPI       *api.Update

	individualRecordingScenario *individualrecording.Client
	webRecordingScenario        *webrecording.Client
	mixRecordingScenario        *mixrecording.Client
}

func NewClient(config *agora.Config) (*Client, error) {
	prefixPath := "/v1/apps/" + config.AppID + "/" + projectName

	c, err := agoraClient.New(config)
	if err != nil {
		return nil, err
	}

	a := &Client{
		acquireAPI:      api.NewAcquire(c, prefixPath),
		startAPI:        api.NewStart("cloudRecording:start", config.Logger, c, prefixPath),
		stopAPI:         api.NewStop(c, prefixPath),
		queryAPI:        api.NewQuery(c, prefixPath),
		updateLayoutAPI: api.NewUpdateLayout(c, prefixPath),
		updateAPI:       api.NewUpdate(c, prefixPath),
	}

	a.individualRecordingScenario = individualrecording.NewClient(a.acquireAPI, a.startAPI, a.stopAPI, a.queryAPI, a.updateLayoutAPI, a.updateAPI)
	a.webRecordingScenario = webrecording.NewClient(a.acquireAPI, a.startAPI, a.stopAPI, a.queryAPI, a.updateLayoutAPI, a.updateAPI)
	a.mixRecordingScenario = mixrecording.NewClient(a.acquireAPI, a.startAPI, a.stopAPI, a.queryAPI, a.updateLayoutAPI, a.updateAPI)
	return a, nil
}

func (a *Client) Acquire(ctx context.Context, payload *api.AcquireReqBody) (*api.AcquireResp, error) {
	return a.acquireAPI.Do(ctx, payload)
}

func (a *Client) Start(ctx context.Context, resourceID string, mode string, payload *api.StartReqBody) (*api.StartResp, error) {
	return a.startAPI.Do(ctx, resourceID, mode, payload)
}

func (a *Client) Stop(ctx context.Context, resourceID string, sid string, mode string, payload *api.StopReqBody) (*api.StopResp, error) {
	return a.stopAPI.Do(ctx, resourceID, sid, mode, payload)
}

func (a *Client) Query(ctx context.Context, resourceID string, sid string, mode string) (*api.QueryResp, error) {
	return a.queryAPI.Do(ctx, resourceID, sid, mode)
}

func (a *Client) Update(ctx context.Context, resourceID string, sid string, mode string, payload *api.UpdateReqBody) (*api.UpdateResp, error) {
	return a.updateAPI.Do(ctx, resourceID, sid, mode, payload)
}

func (a *Client) UpdateLayout() *api.UpdateLayout {
	return a.updateLayoutAPI
}

func (a *Client) IndividualRecording() *individualrecording.Client {
	return a.individualRecordingScenario
}

func (a *Client) WebRecording() *webrecording.Client {
	return a.webRecordingScenario
}

func (a *Client) MixRecording() *mixrecording.Client {
	return a.mixRecordingScenario
}
