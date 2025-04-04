package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
)

type Acquire struct {
	module     string
	logger     log.Logger
	client     client.Client
	prefixPath string // /v1/projects/{appid}/rtsc/cloud-transcoder
}

func NewAcquire(module string, logger log.Logger, client client.Client, prefixPath string) *Acquire {
	return &Acquire{module: module, logger: logger, client: client, prefixPath: prefixPath}
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

	responseData, err := a.doRESTWithRetry(ctx, path, http.MethodPost, payload)
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

const startRetryCount = 3

func (a *Acquire) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		resp       *agora.BaseResponse
		err        error
		retryCount int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		resp, doErr = a.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			a.logger.Debugf(ctx, a.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			)
		default:
			a.logger.Debugf(ctx, a.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= startRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		a.logger.Debugf(ctx, a.module, "http request err:%s", err)
		retryCount++
	})

	return resp, err
}
