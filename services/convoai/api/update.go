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
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Update struct {
	baseHandler
}

func NewUpdate(module string, logger log.Logger, client client.Client, prefixPath string) *Update {
	return &Update{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
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

	responseData, err := u.doRESTWithRetry(ctx, path, http.MethodPost, payload)
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

const updateRetryCount = 3

func (u *Update) doRESTWithRetry(ctx context.Context, path string, method string, requestBody any) (*agora.BaseResponse, error) {
	var (
		baseResponse *agora.BaseResponse
		err          error
		retryCount   int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		baseResponse, doErr = u.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := baseResponse.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			u.logger.Debugf(ctx, u.module, "http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)),
			)
		default:
			u.logger.Debugf(ctx, u.module, "http status code is %d, retry,http response:%s", statusCode, baseResponse.RawBody)
			return fmt.Errorf("http status code is %d, retry", baseResponse.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= updateRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		u.logger.Debugf(ctx, u.module, "http request err:%s", err)
		retryCount++
	})

	return baseResponse, err
}
