package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Delete struct {
	baseHandler
}

func NewDelete(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Delete {
	return &Delete{
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
	responseData, err := doRESTWithRetry(ctx, d.module, d.logger, d.retryCount, d.client, path, http.MethodDelete, nil)
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
