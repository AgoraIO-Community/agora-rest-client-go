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

type Delete struct {
	module     string
	logger     log.Logger
	client     client.Client
	prefixPath string // /v1/projects/{appid}/rtsc/cloud-transcoder
}

func NewDelete(module string, logger log.Logger, client client.Client, prefixPath string) *Delete {
	return &Delete{module: module, logger: logger, client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks/{taskId}?builderToken={tokenName}
func (d *Delete) buildPath(taskId string, tokenName string) string {
	return d.prefixPath + "/tasks/" + taskId + "?builderToken=" + tokenName
}

type DeleteResp struct {
	Response
	SuccessResp DeleteSuccessResp
}

type DeleteSuccessResp struct {
	// Transcoding task ID, a UUID used to identify the cloud transcoder for this request operation
	TaskID string `json:"taskId"`
	
	// Unix timestamp (seconds) when the transcoding task was created
	CreateTs int64 `json:"createTs"`
	
	// Running status of the transcoding task:
	//  - "IDLE": Task has not started.
	//
	//  - "PREPARED": Task has received a start request.
	//
	//  - "STARTING": Task is starting.
	//
	//  - "CREATED": Task initialization complete.
	//
	//  - "STARTED": Task has started.
	//
	//  - "IN_PROGRESS": Task is in progress.
	//
	//  - "STOPPING": Task is stopping.
	//
	//  - "STOPPED": Task has stopped.
	//
	//  - "EXIT": Task exited normally.
	//
	//  - "FAILURE_STOP": Task exited abnormally.
	//
	// Note: You can use this field to monitor the task status.
	Status string `json:"status"`
}

func (d *Delete) Do(ctx context.Context, taskId string, tokenName string) (*DeleteResp, error) {
	path := d.buildPath(taskId, tokenName)
	responseData, err := d.doRESTWithRetry(ctx, path, http.MethodDelete, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp DeleteResp

	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse DeleteSuccessResp
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

const deleteRetryCount = 3

func (d *Delete) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		resp       *agora.BaseResponse
		err        error
		retryCount int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		resp, doErr = d.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			d.logger.Debugf(ctx, d.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			)
		default:
			d.logger.Debugf(ctx, d.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= deleteRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		d.logger.Debugf(ctx, d.module, "http request err:%s", err)
		retryCount++
	})

	return resp, err
}
