package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Query struct {
	baseHandler
}

func NewQuery(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Query {
	return &Query{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			retryCount: retryCount,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}
func (q *Query) buildPath(agentId string) string {
	return q.prefixPath + "/agents/" + agentId
}

func (q *Query) Do(ctx context.Context, agentId string) (*resp.QueryResp, error) {
	path := q.buildPath(agentId)
	responseData, err := doRESTWithRetry(ctx, q.module, q.logger, q.retryCount, q.client, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var queryResp resp.QueryResp

	queryResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.QuerySuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		queryResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		queryResp.ErrResponse = errResponse
	}

	return &queryResp, nil
}
