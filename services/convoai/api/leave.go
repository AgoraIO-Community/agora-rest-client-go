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

type Leave struct {
	baseHandler
}

// NewLeave Creates a new Leave instance
func NewLeave(module string, logger log.Logger, client client.Client, prefixPath string) *Leave {
	return &Leave{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/leave
func (d *Leave) buildPath(agentId string) string {
	return d.prefixPath + "/agents/" + agentId + "/leave"
}

func (d *Leave) Do(ctx context.Context, agentId string) (*resp.LeaveResp, error) {
	path := d.buildPath(agentId)
	responseData, err := d.client.DoREST(ctx, path, http.MethodPost, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var response resp.LeaveResp

	response.BaseResponse = responseData

	if responseData.HttpStatusCode != http.StatusOK {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		response.ErrResponse = errResponse
	}

	return &response, nil
}
