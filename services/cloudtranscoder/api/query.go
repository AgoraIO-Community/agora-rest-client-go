package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Query struct {
	baseHandler
}

func NewQuery(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Query {
	return &Query{
		baseHandler: baseHandler{
			module: module, logger: logger, retryCount: retryCount, client: client, prefixPath: prefixPath,
		},
	}
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

func (q *Query) Do(ctx context.Context, taskId string, tokenName string) (*QueryResp, error) {
	path := q.buildPath(taskId, tokenName)
	responseData, err := doRESTWithRetry(ctx, q.module, q.logger, q.retryCount, q.client, path, http.MethodGet, nil)
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
