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

type Join struct {
	baseHandler
}

func NewJoin(module string, logger log.Logger, client client.Client, prefixPath string) *Join {
	return &Join{
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
// /api/conversational-ai-agent/v2/projects/{appid}/join
func (d *Join) buildPath() string {
	return d.prefixPath + "/join"
}

func (d *Join) Do(ctx context.Context, name string, propertiesBody *req.JoinPropertiesReqBody) (*resp.JoinResp, error) {
	path := d.buildPath()
	request := map[string]any{
		"name":       name,
		"properties": propertiesBody,
	}
	responseData, err := d.client.DoREST(ctx, path, http.MethodPost, request)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var joinResp resp.JoinResp

	joinResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.JoinSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		joinResp.SuccessResp = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		joinResp.ErrResponse = errResponse
	}

	return &joinResp, nil
}
