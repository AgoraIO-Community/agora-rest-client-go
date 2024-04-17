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
	// UID 字符串内容为用户 UID
	UID string `json:"uid"`

	// ImageURL 字符串内容该用户背景图片的 URL 地址
	//
	// 配置背景图后，当该⽤户停止发送视频流超过 3.5 秒，画⾯将切换为该背景图。
	//
	//URL 支持 HTTP 和 HTTPS 协议，图片格式支持 JPG 和 BMP。图片大小不得超过 6 MB。录制服务成功下载图片后，设置才会生效；如果下载失败，则设置不⽣效。不同字段设置可能会互相覆盖，具体规则详见设置背景色或背景图。
	ImageURL string `json:"image_url"`

	// RenderMode 画面显示模式
	//
	//0：裁剪模式。优先保证画面被填满。视频尺寸等比缩放，直至整个画面被视频填满。如果视频长宽与显示窗口不同，则视频流会按照画面设置的比例进行周边裁剪后填满画面。
	//
	//1：缩放模式。优先保证视频内容全部显示。视频尺寸等比缩放，直至视频窗口的一边与画面边框对齐。如果视频尺寸与画面尺寸不一致，在保持长宽比的前提下，将视频进行缩放后填满画面，缩放后的视频四周会有一圈黑边。
	//
	// 默认值：0
	RenderMode int `json:"render_mode"`
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
