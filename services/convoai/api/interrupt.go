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

type Interrupt struct {
	baseHandler
}

func NewInterrupt(module string, logger log.Logger, client client.Client, prefixPath string) *Interrupt {
	return &Interrupt{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
//
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/interrupt
func (i *Interrupt) buildPath(agentId string) string {
	return i.prefixPath + "/agents/" + agentId + "/interrupt"
}

func (i *Interrupt) Do(ctx context.Context, agentId string) (*resp.InterruptResp, error) {
	path := i.buildPath(agentId)
	responseData, err := i.client.DoREST(ctx, path, http.MethodPost, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var interruptResp resp.InterruptResp

	interruptResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.InterruptSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		interruptResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		interruptResp.ErrResponse = errResponse
	}

	return &interruptResp, nil
}
