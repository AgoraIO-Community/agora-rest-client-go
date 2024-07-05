package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"

	"github.com/tidwall/gjson"
)

type Acquire struct {
	client     core.Client
	prefixPath string // /v1/projects/{appid}/rtsc/cloud-transcoder
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/buildTokens
func (a *Acquire) buildPath() string {
	return a.prefixPath + "/buildTokens"
}

type AcquireReqBody struct {
	// 用户指定的实例 ID。长度必须在 64 个字符以内，支持的字符集范围为：
	//
	//   - 所有小写英文字母（a-z）
	//
	//   - 所有大写英文字母（A-Z）
	//
	//   - 数字 0-9
	//
	//   - "-", "_"
	//
	// 注意：一个 instanceId 可以生成多个 builderToken，但在一个任务中只能使用一个 builderToken 发起请求。
	InstanceId string `json:"instanceId"`
}

type AcquireResp struct {
	Response
	SuccessResp AcquireSuccessResp
}

type AcquireSuccessResp struct {
	// 生成 builderToken 时的 Unix 时间戳（秒）
	CreateTs string `json:"createTs"`
	// 请求时设置的 instanceId
	InstanceId string `json:"instanceId"`
	// 代表 builderToken 的值，在后续调用其他方法时需要传入该值
	TokenName string `json:"tokenName"`
}

func (a *Acquire) Do(ctx context.Context, payload *AcquireReqBody) (*AcquireResp, error) {
	path := a.buildPath()

	responseData, err := a.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp AcquireResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse AcquireSuccessResp
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
