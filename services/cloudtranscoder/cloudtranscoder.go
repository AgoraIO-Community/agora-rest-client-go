package cloudtranscoder

import (
	"context"

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

func NewClient(config *agora.Config) (*Client, error) {
	prefixPath := "/v1/projects/" + config.AppID + "/" + projectName
	c, err := agoraClient.New(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		acquireAPI: api.NewAcquire("cloudTranscoder:acquire", config.Logger, c, prefixPath),
		createAPI:  api.NewCreate("cloudTranscoder:create", config.Logger, c, prefixPath),
		queryAPI:   api.NewQuery("cloudTranscoder:query", config.Logger, c, prefixPath),
		deleteAPI:  api.NewDelete("cloudTranscoder:delete", config.Logger, c, prefixPath),
		updateAPI:  api.NewUpdate("cloudTranscoder:update", config.Logger, c, prefixPath),
	}, nil
}

func (a *Client) Acquire(ctx context.Context, payload *api.AcquireReqBody) (*api.AcquireResp, error) {
	return a.acquireAPI.Do(ctx, payload)
}

func (a *Client) Create(ctx context.Context, tokenName string, payload *api.CreateReqBody) (*api.CreateResp, error) {
	return a.createAPI.Do(ctx, tokenName, payload)
}

func (a *Client) Query(ctx context.Context, taskId string, tokenName string) (*api.QueryResp, error) {
	return a.queryAPI.Do(ctx, taskId, tokenName)
}

func (a *Client) Delete(ctx context.Context, taskId string, tokenName string) (*api.DeleteResp, error) {
	return a.deleteAPI.Do(ctx, taskId, tokenName)
}

func (a *Client) Update(ctx context.Context, taskId string, tokenName string, sequenceId uint, updateMask string,
	payload *api.UpdateReqBody) (*api.UpdateResp, error) {
	return a.updateAPI.Do(ctx, taskId, tokenName, sequenceId, updateMask, payload)
}
