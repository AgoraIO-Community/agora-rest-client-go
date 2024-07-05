package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"

	"github.com/tidwall/gjson"
)

type Acquire struct {
	client     core.Client
	prefixPath string // /v1/projects/{appid}/rtsc/cloud-transcoder
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/buildTokens
func (a *Acquire) buildPath() string {
	return a.prefixPath + "/buildTokens"
}

type AcquireReqBody struct {
	InstanceId string `json:"instanceId"`
}

type AcquireResp struct {
	Response
	SuccessResp AcquireSuccessResp
}

type AcquireSuccessResp struct {
	CreateTs   string `json:"createTs"`
	InstanceId string `json:"instanceId"`
	TokenName  string `json:"tokenName"`
}

func (a *Acquire) Do(ctx context.Context, payload *AcquireReqBody) (*AcquireResp, error) {
	path := a.buildPath()

	responseData, err := a.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp AcquireResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse AcquireSuccessResp
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
