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
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Query struct {
	baseHandler
}

func NewQuery(module string, logger log.Logger, client client.Client, prefixPath string) *Query {
	return &Query{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}
func (q *Query) buildPath(agentId string) string {
	return q.prefixPath + "/agents/" + agentId
}

func (q *Query) Do(ctx context.Context, agentId string) (*resp.QueryResp, error) {
	path := q.buildPath(agentId)
	responseData, err := q.doRESTWithRetry(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var queryResp resp.QueryResp

	queryResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.QuerySuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		queryResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		queryResp.ErrResponse = errResponse
	}

	return &queryResp, nil
}

const queryRetryCount = 3

func (q *Query) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		baseResponse *agora.BaseResponse
		err          error
		retryCount   int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		baseResponse, doErr = q.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := baseResponse.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			q.logger.Debugf(ctx, q.module, "http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)),
			)
		default:
			q.logger.Debugf(ctx, q.module, "http status code is %d, retry,http response:%s", statusCode, baseResponse.RawBody)
			return fmt.Errorf("http status code is %d, retry", baseResponse.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= queryRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		q.logger.Debugf(ctx, q.module, "http request err:%s", err)
		retryCount++
	})

	return baseResponse, err
}
