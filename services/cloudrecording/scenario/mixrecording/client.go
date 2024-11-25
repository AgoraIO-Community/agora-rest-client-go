package mixrecording

import (
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

func (i *Client) Acquire() *Acquire {
	return NewAcquire(i.acquireAPI)
}

func (i *Client) Query() *Query {
	return NewQuery(i.queryAPI)
}

func (i *Client) Start() *Start {
	return NewStart(i.startAPI)
}

func (i *Client) Stop() *Stop {
	return NewStop(i.stopAPI)
}

func (i *Client) Update() *Update {
	return NewUpdate(i.updateAPI)
}

func (i *Client) UpdateLayout() *UpdateLayout {
	return NewUpdateLayout(i.updateLayoutAPI)
}
