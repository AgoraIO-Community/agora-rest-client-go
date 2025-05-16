package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Acquire struct {
	baseHandler
}

func NewAcquire(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Acquire {
	return &Acquire{
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
// /v1/projects/{appid}/rtsc/cloud-transcoder/builderTokens
func (a *Acquire) buildPath() string {
	return a.prefixPath + "/builderTokens"
}

type AcquireReqBody struct {
	// User-specified instance ID. Must be within 64 characters, supporting:
	//   - All lowercase letters (a-z)
	//   - All uppercase letters (A-Z)
	//   - Numbers 0-9
	//   - "-", "_"
	// Note: One instanceId can generate multiple builderTokens, but only one can be used per task.
	InstanceId string             `json:"instanceId,omitempty"`
	Services   *CreateReqServices `json:"services,omitempty"`
}

type AcquireResp struct {
	Response
	SuccessResp AcquireSuccessResp
}

type AcquireSuccessResp struct {
	// Unix timestamp (seconds) when the builderToken was generated
	CreateTs int64 `json:"createTs"`
	// The instanceId set in the request
	InstanceId string `json:"instanceId"`
	// The value of the builderToken, required for subsequent method calls
	TokenName string `json:"tokenName"`
}

func (a *Acquire) Do(ctx context.Context, payload *AcquireReqBody) (*AcquireResp, error) {
	path := a.buildPath()

	responseData, err := doRESTWithRetry(ctx, a.module, a.logger, a.retryCount, a.client, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp AcquireResp

	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse AcquireSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResponse
	} else {
		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	return &resp, nil
}
