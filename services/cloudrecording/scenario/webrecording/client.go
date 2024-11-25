package webrecording

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

func (c *Client) Acquire() *Acquire {
	return NewAcquire(c.acquireAPI)
}

func (c *Client) Query() *Query {
	return NewQuery(c.queryAPI)
}

func (c *Client) Start() *Start {
	return NewStart(c.startAPI)
}

func (c *Client) Stop() *Stop {
	return NewStop(c.stopAPI)
}

func (c *Client) Update() *Update {
	return NewUpdate(c.updateAPI)
}
