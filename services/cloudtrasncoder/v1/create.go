package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/tidwall/gjson"
)

type Create struct {
	client     core.Client
	prefixPath string
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks
func (c *Create) buildPath(tokenName string) string {
	return c.prefixPath + "/tasks?builderToken=" + tokenName
}

type CreateReqBody struct {
	Services CreateReqServices `json:"services"`
}

type CreateReqServices struct {
	CloudTrasncoder *CloudTrasncoderPayload `json:"cloudTranscoder"`
}

type CloudTrasncoderPayload struct {
	ServiceType string `json:"serviceType"`
}

type CreateSuccessResp struct {
	TaskID   string `json:"taskId"`
	CreateTs int64  `json:"createTs"`
	Status   string `json:"status"`
}

type CreateResp struct {
	Response
	SuccessResp CreateSuccessResp
}

func (c *Create) Do(ctx context.Context, tokenName string, payload *CreateReqBody) (*CreateResp, error) {
	path := c.buildPath(tokenName)
	responseData, err := c.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp CreateResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse CreateSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResponse
	} else {
		codeResult := gjson.GetBytes(responseData.RawBody, "code")
		if !codeResult.Exists() {
			return nil, core.NewGatewayErr(responseData.HttpStatusCode, string(responseData.RawBody))
		}
		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	resp.BaseResponse = responseData

	return &resp, nil
}
