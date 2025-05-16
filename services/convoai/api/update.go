package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Update struct {
	baseHandler
}

func NewUpdate(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Update {
	return &Update{
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
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/update
func (u *Update) buildPath(agentId string) string {
	return u.prefixPath + "/agents/" + agentId + "/update"
}

func (u *Update) Do(ctx context.Context, agentId string, payload *req.UpdateReqBody) (*resp.UpdateResp, error) {
	path := u.buildPath(agentId)

	responseData, err := doRESTWithRetry(ctx, u.module, u.logger, u.retryCount, u.client, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var updateResp resp.UpdateResp

	updateResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.UpdateSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		updateResp.SuccessResp = successResponse
		return &updateResp, nil
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		updateResp.ErrResponse = errResponse
	}

	return &updateResp, nil
}
