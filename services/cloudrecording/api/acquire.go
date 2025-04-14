package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type Acquire struct {
	client     client.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording/
}

func NewAcquire(client client.Client, prefixPath string) *Acquire {
	return &Acquire{client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/acquire
func (a *Acquire) buildPath() string {
	return a.prefixPath + "/acquire"
}

type AcquireReqBody struct {
	Cname         string                `json:"cname"`
	Uid           string                `json:"uid"`
	ClientRequest *AcquireClientRequest `json:"clientRequest"`
}

type AcquireClientRequest struct {
	Scene int `json:"scene"`

	StartParameter      *StartClientRequest `json:"startParameter,omitempty"`
	ResourceExpiredHour int                 `json:"resourceExpiredHour,omitempty"`
	ExcludeResourceIds  []string            `json:"excludeResourceIds,omitempty"`
	RegionAffinity      int                 `json:"regionAffinity,omitempty"`
}

// @brief AcquireResp returned by the various of cloud recording scenarios Acquire API.
//
// @since v0.8.0
type AcquireResp struct {
	// Response returned by the cloud recording API, see Response for details
	Response
	// Successful response, see AcquireSuccessResp for details
	SuccessRes AcquireSuccessResp
}

// @brief Successful response returned by the various of cloud recording scenarios Acquire API.
//
// @since v0.8.0
type AcquireSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string `json:"resourceId"`
}

func (a *Acquire) Do(ctx context.Context, payload *AcquireReqBody) (*AcquireResp, error) {
	path := a.buildPath()

	responseData, err := a.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
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
		resp.SuccessRes = successResponse
	} else {
		codeResult := gjson.GetBytes(responseData.RawBody, "code")
		if !codeResult.Exists() {
			return nil, agora.NewGatewayErr(responseData.HttpStatusCode, string(responseData.RawBody))
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
