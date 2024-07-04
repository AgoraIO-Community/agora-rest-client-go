package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/tidwall/gjson"
)

type Update struct {
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
	SuccessResp UpdateSuccessResp
}

type UpdateSuccessResp struct{}

func (a *Update) Do(ctx context.Context, taskId string, tokenName string, sequenceId uint, payload *UpdateReqBody) (*UpdateResp, error) {
	path := a.buildPath(taskId, tokenName, sequenceId)

	responseData, err := a.client.DoREST(ctx, path, http.MethodPatch, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp UpdateResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse UpdateSuccessResp
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
