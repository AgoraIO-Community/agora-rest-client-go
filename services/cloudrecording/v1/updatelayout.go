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
	// MaxResolutionUID 仅需在垂直布局下设置。指定显示大视窗画面的用户 UID。
	//
	// 字符串内容的整型取值范围 1 到 (2^32-1)，且不可设置为 0。
	MaxResolutionUID string `json:"maxResolutionUid,omitempty"`

	// MixedVideoLayout 视频合流布局：
	//
	// 0：悬浮布局。第一个加入频道的用户在屏幕上会显示为大视窗，铺满整个画布，其他用户的视频画面会显示为小视窗，从下到上水平排列，最多 4 行，每行 4 个画面，最多支持共 17 个画面。
	//
	// 1：自适应布局。根据用户的数量自动调整每个画面的大小，每个用户的画面大小一致，最多支持 17 个画面。
	//
	// 2：垂直布局。指定 maxResolutionUid 在屏幕左侧显示大视窗画面，其他用户的小视窗画面在右侧垂直排列，最多两列，一列 8 个画面，最多支持共 17 个画面。
	//
	// 3：自定义布局。由你在 layoutConfig 字段中自定义合流布局。
	MixedVideoLayout int `json:"mixedVideoLayout"`

	// BackgroundColor 视频画布的背景颜色。
	//
	// 支持 RGB 颜色表，字符串格式为 # 号和 6 个十六进制数。
	//
	// 默认值: #000000
	BackgroundColor string `json:"backgroundColor,omitempty"`

	// BackgroundImage 视频画布的背景图的 URL。背景图的显示模式为裁剪模式。
	//
	// 裁剪模式：优先保证画面被填满。背景图尺寸等比缩放，直至整个画面被背景图填满。如果背景图长宽与显示窗口不同，则背景图会按照画面设置的比例进行周边裁剪后填满画面。
	BackgroundImage string `json:"backgroundImage,omitempty"`

	// DefaultUserBackgroundImage 默认的用户画面背景图的 URL。
	//
	//配置该字段后，当任一⽤户停止发送视频流超过 3.5 秒，画⾯将切换为该背景图；如果针对某 UID 单独设置了背景图，则该设置会被覆盖。
	DefaultUserBackgroundImage string `json:"defaultUserBackgroundImage,omitempty"`

	LayoutConfig     []UpdateLayoutConfig `json:"layoutConfig,omitempty"`
	BackgroundConfig []BackgroundConfig   `json:"backgroundConfig,omitempty"`
}

// UpdateLayoutConfig 用户的合流画面布局。由每个用户对应的布局画面设置组成的数组，支持最多 17 个用户。
type UpdateLayoutConfig struct {
	// UID 字符串内容为待显示在该区域的用户的 UID，32 位无符号整数。
	//
	// 如果不指定 UID，会按照用户加入频道的顺序自动匹配 layoutConfig 中的画面设置。
	UID string `json:"uid"`

	// XAxis 屏幕里该画面左上角的横坐标的相对值，精确到小数点后六位。从左到右布局，0.0 在最左端，1.0 在最右端。
	//
	// 该字段也可以设置为整数 0 或 1。
	//
	// 0 <= XAxis <= 1
	XAxis float32 `json:"x_axis"`

	// YAxis 屏幕里该画面左上角的纵坐标的相对值，精确到小数点后六位。从上到下布局，0.0 在最上端，1.0 在最下端。
	//
	// 该字段也可以设置为整数 0 或 1。
	//
	// 0 <= YAxis <= 1
	YAxis float32 `json:"y_axis"`

	// Width 该画面宽度的相对值，精确到小数点后六位。该字段也可以设置为整数 0 或 1。
	//
	// 0 <= Width <= 1
	Width float32 `json:"width"`

	// Height 该画面高度的相对值，精确到小数点后六位。该字段也可以设置为整数 0 或 1。
	//
	// 0 <= Height <= 1
	Height float32 `json:"height"`

	// Alpha 图像的透明度。精确到小数点后六位。
	//
	// 0.0 表示图像为透明的，1.0 表示图像为完全不透明的。
	//
	// 0 <= Alpha <= 1
	// 默认值：1
	Alpha float32 `json:"alpha"`

	// RenderMode 画面显示模式
	//
	// 0：裁剪模式。优先保证画面被填满。视频尺寸等比缩放，直至整个画面被视频填满。如果视频长宽与显示窗口不同，则视频流会按照画面设置的比例进行周边裁剪后填满画面。
	//
	// 1：缩放模式。优先保证视频内容全部显示。视频尺寸等比缩放，直至视频窗口的一边与画面边框对齐。如果视频尺寸与画面尺寸不一致，在保持长宽比的前提下，将视频进行缩放后填满画面，缩放后的视频四周会有一圈黑边。
	//
	// 默认值：0
	RenderMode int `json:"render_mode"`
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
