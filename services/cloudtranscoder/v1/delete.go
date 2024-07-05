package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"

	"github.com/tidwall/gjson"
)

type Delete struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
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

func (d *Delete) Do(ctx context.Context, taskId string, tokenName string) (*DeleteResp, error) {
	path := d.buildPath(taskId, tokenName)
	responseData, err := d.client.DoREST(ctx, path, http.MethodDelete, nil)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp DeleteResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse DeleteSuccessResp
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
