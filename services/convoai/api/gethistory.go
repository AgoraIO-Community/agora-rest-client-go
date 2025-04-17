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

type History struct {
	baseHandler
}

func NewHistory(module string, logger log.Logger, client client.Client, prefixPath string) *History {
	return &History{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/history
func (h *History) buildPath(agentId string) string {
	return h.prefixPath + "/agents/" + agentId + "/history"
}

func (h *History) Do(ctx context.Context, agentId string) (*resp.HistoryResp, error) {
	path := h.buildPath(agentId)
	responseData, err := h.client.DoREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var historyResp resp.HistoryResp

	historyResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.HistorySuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		historyResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		historyResp.ErrResponse = errResponse
	}

	return &historyResp, nil
}
