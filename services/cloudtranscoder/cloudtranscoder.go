package cloudtranscoder

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	agoraClient "github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudtranscoder/api"
)

const projectName = "rtsc/cloud-transcoder"

type Client struct {
	acquireAPI *api.Acquire
	createAPI  *api.Create
	queryAPI   *api.Query
	deleteAPI  *api.Delete
	updateAPI  *api.Update
}

func NewClient(config *agora.Config) *Client {
	prefixPath := "/v1/projects/" + config.AppID + "/" + projectName
	c := agoraClient.New(config)

	return &Client{
		acquireAPI: api.NewAcquire("cloudTranscoder:acquire", config.Logger, c, prefixPath),
		createAPI:  api.NewCreate("cloudTranscoder:create", config.Logger, c, prefixPath),
		queryAPI:   api.NewQuery("cloudTranscoder:query", config.Logger, c, prefixPath),
		deleteAPI:  api.NewDelete("cloudTranscoder:delete", config.Logger, c, prefixPath),
		updateAPI:  api.NewUpdate("cloudTranscoder:update", config.Logger, c, prefixPath),
	}
}

func (a *Client) Acquire() *api.Acquire {
	return a.acquireAPI
}

func (a *Client) Create() *api.Create {
	return a.createAPI
}

func (a *Client) Query() *api.Query {
	return a.queryAPI
}

func (a *Client) Delete() *api.Delete {
	return a.deleteAPI
}

func (a *Client) Update() *api.Update {
	return a.updateAPI
}
