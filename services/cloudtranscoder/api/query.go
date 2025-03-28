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

type Query struct {
	module     string
	logger     log.Logger
	client     client.Client
	prefixPath string // /v1/projects/{appid}/rtsc/cloud-transcoder
}

func NewQuery(module string, logger log.Logger, client client.Client, prefixPath string) *Query {
	return &Query{module: module, logger: logger, client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks/{taskId}?builderToken={tokenName}
func (q *Query) buildPath(taskId string, tokenName string) string {
	return q.prefixPath + "/tasks/" + taskId + "?builderToken=" + tokenName
}

type QueryResp struct {
	Response
	SuccessRes QuerySuccessResp
}

type QuerySuccessResp struct {
	// 转码任务 ID，为 UUID，用于标识本次请求操作的 cloud transcoder
	TaskID string `json:"taskId"`
	// 转码任务创建时的 Unix 时间戳（秒）
	CreateTs int64 `json:"createTs"`
	// 转码任务的运行状态：
	//  - "IDLE": 任务未开始。
	//
	//  - "PREPARED": 任务已收到开启请求。
	//
	//  - "STARTING": 任务正在开启。
	//
	//  - "CREATED": 任务初始化完成。
	//
	//  - "STARTED": 任务已经启动。
	//
	//  - "IN_PROGRESS": 任务正在进行。
	//
	//  - "STOPPING": 任务正在停止。
	//
	//  - "STOPPED": 任务已经停止。
	//
	//  - "EXIT": 任务正常退出。
	//
	//  - "FAILURE_STOP": 任务异常退出。
	//
	// 注意：你可以用该字段监听任务的状态。
	Status string `json:"status"`
}

func (q *Query) Do(ctx context.Context, taskId string, tokenName string) (*QueryResp, error) {
	path := q.buildPath(taskId, tokenName)
	responseData, err := q.doRESTWithRetry(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp QueryResp

	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse QuerySuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessRes = successResponse
	} else {
		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	return &resp, nil
}

const queryRetryCount = 3

func (q *Query) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		resp       *agora.BaseResponse
		err        error
		retryCount int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		resp, doErr = q.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			q.logger.Debugf(ctx, q.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			)
		default:
			q.logger.Debugf(ctx, q.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
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

	return resp, err
}
