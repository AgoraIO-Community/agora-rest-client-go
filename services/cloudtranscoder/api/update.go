package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Update struct {
	baseHandler
}

func NewUpdate(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Update {
	return &Update{
		baseHandler: baseHandler{
			module: module, logger: logger, retryCount: retryCount, client: client, prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks/{taskId}?token_name={tokenName}&sequence_id={sequenceId}&updateMask=services.cloudTranscoder.config
func (u *Update) buildPath(taskId string, tokenName string, sequenceId uint, updateMask string) string {
	return u.prefixPath + "/tasks/" + taskId + "?" +
		"builderToken=" + tokenName +
		"&sequenceId=" + strconv.Itoa(int(sequenceId)) +
		"&updateMask=" + updateMask
}

type UpdateReqBody struct {
	Services *CreateReqServices `json:"services"`
}

type UpdateResp struct {
	Response
}

func (u *Update) Do(ctx context.Context, taskId string, tokenName string, sequenceId uint, updateMask string, payload *UpdateReqBody) (*UpdateResp, error) {
	path := u.buildPath(taskId, tokenName, sequenceId, updateMask)

	responseData, err := doRESTWithRetry(ctx, u.module, u.logger, u.retryCount, u.client, path, http.MethodPatch, payload)
	if err != nil {
		var internalErr *agora.InternalErr
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
