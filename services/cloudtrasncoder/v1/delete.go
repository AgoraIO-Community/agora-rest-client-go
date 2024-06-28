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
	TaskID   string `json:"taskId"`
	CreateTs int64  `json:"createTs"`
	Status   string `json:"status"`
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
