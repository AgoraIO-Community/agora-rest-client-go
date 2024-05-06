package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Acquire struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording/
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/acquire
func (a *Acquire) buildPath() string {
	return a.prefixPath + "/acquire"
}

type AcquirerReqBody struct {
	Cname         string                 `json:"cname"`
	Uid           string                 `json:"uid"`
	ClientRequest *AcquirerClientRequest `json:"clientRequest"`
}

type AcquirerClientRequest struct {
	Scene int `json:"scene"`

	// StartParameter 设置该字段后，可以提升可用性并优化负载均衡。
	//
	// 注意：如果填写该字段，则必须确保 startParameter object 和后续 start 请求中填写的 clientRequest object 完全一致，
	// 且取值合法，否则 start 请求会收到报错。
	StartParameter      *StartClientRequest `json:"startParameter,omitempty"`
	ResourceExpiredHour int                 `json:"resourceExpiredHour"`
	ExcludeResourceIds  []string            `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int                 `json:"regionAffinity,omitempty"`
}

type AcquirerResp struct {
	Response
	SuccessRes AcquirerSuccessResp
}

type AcquirerSuccessResp struct {
	ResourceId string `json:"resourceId"`
}

func (a *Acquire) Do(ctx context.Context, payload *AcquirerReqBody) (*AcquirerResp, error) {
	path := a.buildPath()

	responseData, err := a.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp AcquirerResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse AcquirerSuccessResp
		if err = responseData.UnmarshallToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessRes = successResponse
	} else {
		codeResult := gjson.GetBytes(responseData.RawBody, "code")
		if !codeResult.Exists() {
			return nil, core.NewGatewayErr(responseData.HttpStatusCode, string(responseData.RawBody))
		}
		var errResponse ErrResponse
		if err = responseData.UnmarshallToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	resp.BaseResponse = responseData
	return &resp, nil
}
