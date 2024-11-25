package cloudrecording

import (
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

func NewClient(config *agora.Config) *Client {
	prefixPath := "/v1/apps/" + config.AppID + "/" + projectName
	c := agoraClient.New(config)
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
	return a
}

func (a *Client) Acquire() *api.Acquire {
	return a.acquireAPI
}

func (a *Client) Start() *api.Start {
	return a.startAPI
}

func (a *Client) Stop() *api.Stop {
	return a.stopAPI
}

func (a *Client) Query() *api.Query {
	return a.queryAPI
}

func (a *Client) Update() *api.Update {
	return a.updateAPI
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
