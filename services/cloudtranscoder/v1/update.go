package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Update struct {
	module     string
	logger     core.Logger
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks/{taskId}?token_name={tokenName}&sequence_id={sequenceId}&updateMask=services.cloudTranscoder.config
func (u *Update) buildPath(taskId string, tokenName string, sequenceId uint) string {
	return u.prefixPath + "/tasks/" + taskId + "?" +
		"builderToken=" + tokenName +
		"&sequenceId=" + strconv.Itoa(int(sequenceId)) +
		"&updateMask=services.cloudTranscoder.config"
}

type UpdateReqBody struct {
	Services CreateReqServices `json:"services"`
}

type UpdateResp struct {
	Response
}

func (u *Update) Do(ctx context.Context, taskId string, tokenName string, sequenceId uint, payload *UpdateReqBody) (*UpdateResp, error) {
	path := u.buildPath(taskId, tokenName, sequenceId)

	responseData, err := u.doRESTWithRetry(ctx, path, http.MethodPatch, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp UpdateResp

	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		return &resp, nil
	} else {
		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	return &resp, nil
}

const updateRetryCount = 3

func (u *Update) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*core.BaseResponse, error) {
	var (
		resp  *core.BaseResponse
		err   error
		retry int
	)

	err = core.RetryDo(func(retryCount int) error {
		var doErr error

		resp, doErr = u.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return core.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			u.logger.Debugf(ctx, u.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return core.NewRetryErr(
				false,
				core.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			)
		default:
			u.logger.Debugf(ctx, u.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retry >= updateRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		u.logger.Debugf(ctx, u.module, "http request err:%s", err)
		retry++
	})

	return resp, err
}
