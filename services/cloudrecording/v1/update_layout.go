package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type UpdateLayout struct {
	client     core.Client
	prefixPath string // /apps/{appid}/cloud_recording
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/updateLayout
func (s *UpdateLayout) buildPath(resourceID string, sid string, mode string) string {
	return s.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/updateLayout"
}

type UpdateLayoutReqBody struct {
	Cname         string                     `json:"cname"`
	Uid           string                     `json:"uid"`
	ClientRequest *UpdateLayoutClientRequest `json:"clientRequest"`
}

type UpdateLayoutClientRequest struct {
	MaxResolutionUID           string               `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout           int                  `json:"mixedVideoLayout"`
	BackgroundColor            string               `json:"backgroundColor,omitempty"`
	BackgroundImage            string               `json:"backgroundImage,omitempty"`
	DefaultUserBackgroundImage string               `json:"defaultUserBackgroundImage,omitempty"`
	LayoutConfig               []UpdateLayoutConfig `json:"layoutConfig,omitempty"`
	BackgroundConfig           []BackgroundConfig   `json:"backgroundConfig,omitempty"`
}

type UpdateLayoutConfig struct {
	UID        string  `json:"uid"`
	XAxis      float32 `json:"x_axis"`
	YAxis      float32 `json:"y_axis"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Alpha      float32 `json:"alpha"`
	RenderMode int     `json:"render_mode"`
}

type BackgroundConfig struct {
	UID        string `json:"uid"`
	ImageURL   string `json:"image_url"`
	RenderMode int    `json:"render_mode"`
}

type UpdateLayoutSuccessResp struct {
	ResourceId string `json:"resourceId"`
	SID        string `json:"sid"`
}

type UpdateLayoutResp struct {
	Response
	SuccessResp UpdateLayoutSuccessResp
}

func (s *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateLayoutReqBody) (*UpdateLayoutResp, error) {
	path := s.buildPath(resourceID, sid, mode)

	responseData, err := s.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp UpdateLayoutResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse UpdateLayoutSuccessResp
		if err = responseData.UnmarshallToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResponse
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
