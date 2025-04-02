package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type UpdateLayout struct {
	client     client.Client
	prefixPath string // /apps/{appid}/cloud_recording
}

func NewUpdateLayout(client client.Client, prefixPath string) *UpdateLayout {
	return &UpdateLayout{client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/updateLayout
func (u *UpdateLayout) buildPath(resourceID string, sid string, mode string) string {
	return u.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/updateLayout"
}

type UpdateLayoutReqBody struct {
	Cname         string                     `json:"cname"`
	Uid           string                     `json:"uid"`
	ClientRequest *UpdateLayoutClientRequest `json:"clientRequest"`
}

type UpdateLayoutClientRequest struct {
	MaxResolutionUID          string               `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout         int                  `json:"mixedVideoLayout"`
	BackgroundColor         string               `json:"backgroundColor,omitempty"`
	BackgroundImage         string               `json:"backgroundImage,omitempty"`
	DefaultUserBackgroundImage string               `json:"defaultUserBackgroundImage,omitempty"`
	LayoutConfig            []UpdateLayoutConfig `json:"layoutConfig,omitempty"`
	BackgroundConfig        []BackgroundConfig   `json:"backgroundConfig,omitempty"`
}

type UpdateLayoutConfig struct {
	UID        string  `json:"uid"`

	// XAxis is the relative x-coordinate of the top-left corner of the video.
	// 0 <= XAxis <= 1
	XAxis      float32 `json:"x_axis"`

	// YAxis is the relative y-coordinate of the top-left corner of the video.
	// 0 <= YAxis <= 1
	YAxis      float32 `json:"y_axis"`

	// Width is the relative width of the video.
	// 0 <= Width <= 1
	Width      float32 `json:"width"`

	// Height is the relative height of the video.
	// 0 <= Height <= 1
	Height     float32 `json:"height"`

	// Alpha is the transparency of the video image.
	// 0 <= Alpha <= 1
	Alpha      float32 `json:"alpha"`

	// RenderMode specifies the video display mode
	// 0: Cropping mode
	// 1: Scaling mode
	RenderMode int `json:"render_mode"`
}

type BackgroundConfig struct {
	UID        string `json:"uid"`
	ImageURL   string `json:"image_url"`
	RenderMode int    `json:"render_mode"`
}

type UpdateLayoutSuccessResp struct {
	ResourceId string `json:"resourceId"`
	Sid        string `json:"sid"`
}

type UpdateLayoutResp struct {
	Response
	SuccessResponse UpdateLayoutSuccessResp
}

func (u *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateLayoutReqBody) (*UpdateLayoutResp, error) {
	path := u.buildPath(resourceID, sid, mode)

	responseData, err := u.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp UpdateLayoutResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse UpdateLayoutSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResponse = successResponse
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
