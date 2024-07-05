package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"

	"github.com/tidwall/gjson"
)

type Create struct {
	client     core.Client
	prefixPath string
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks?builderToken={tokenName}
func (c *Create) buildPath(tokenName string) string {
	return c.prefixPath + "/tasks?builderToken=" + tokenName
}

type CreateReqBody struct {
	Services CreateReqServices `json:"services"`
}

type CreateReqServices struct {
	CloudTranscoder *CloudTranscoderPayload `json:"cloudTranscoder"`
}

type CloudTranscoderPayload struct {
	// 服务类型，此处为 "cloudTranscoderV2"
	ServiceType string                 `json:"serviceType"`
	Config      *CloudTranscoderConfig `json:"config"`
}

type CloudTranscoderConfig struct {
	Transcoder *CloudTranscoderConfigPayload `json:"transcoder"`
}

type CloudTranscoderConfigPayload struct {
	// Cloud transcoder 处于空闲状态的最大时长（秒）。空闲指 cloud transcoder 处理的音视频流所对应的所有主播均已离开频道。
	// 空闲状态超过设置的 idleTimeOut 后， cloud transcoder 会自动销毁。
	//
	// 范围：[1,86400]
	//
	// 默认值:300
	IdleTimeout uint                        `json:"idleTimeout"`
	AudioInputs []CloudTranscoderAudioInput `json:"audioInputs"`
	VideoInputs []CloudTranscoderVideoInput `json:"videoInputs"`
	Canvas      *CloudTranscoderCanvas      `json:"canvas"`
	WaterMarks  []CloudTranscoderWaterMark  `json:"waterMarks"`
	Outputs     []CloudTranscoderOutput     `json:"outputs"`
}

type CloudTranscoderAudioInput struct {
	Rtc *CloudTranscoderRtc `json:"rtc"`
}

type CloudTranscoderRtc struct {
	// 音视频输入源（或输出）所属的 RTC 频道名
	//
	// 目前仅支持订阅单个频道的音视频源，音频源和视频源所属频道必须相同。
	RtcChannel string `json:"rtcChannel"`
	// 音视频输入源（或输出）所对应的 UID
	//
	// RTC 频道内不允许存在相同的 UID，因此，请确保该值与频道内其他用户 UID 不同。
	RtcUID int `json:"rtcUid"`
	// Cloud transcoder 在进入待转码视频源（或转码输出音视频流）所属 RTC 频道时所需设置的 Token。
	//
	// 该值可用于确保频道安全，避免异常用户扰乱频道内其他用户。
	//
	// 注意：
	//   - 当前配置输入流的时，Cloud transcoder 在待转码音视频源所属 RTC 频道内的 UID 为声网随机分配。因此，生成 Token 时，你使用的 uid 必须为 0。
	//   - 当前配置输出流的时，Cloud transcoder 在转码输出音视频流所属 RTC 频道内的 UID 为你指定的 outputs.rtc.rtcUid，因此，生成 Token 时，你使用的 uid 必须和 outputs.rtc.rtcUid 一致
	RtcToken string `json:"rtcToken"`
}

type CloudTranscoderVideoInput struct {
	Rtc *CloudTranscoderRtc `json:"rtc"`
	// 用户离线时占位图片的 URL。
	//
	// 必须为合法 URL，且包含 jpg 或 png 后缀。
	PlaceholderImageURL string                 `json:"placeholderImageUrl"`
	Region              *CloudTranscoderRegion `json:"region"`
}

type CloudTranscoderRegion struct {
	// 画面在画布上的 x 坐标 (px)。
	//
	// 以画布左上角为原点，x 坐标为画面左上角相对于原点的横向位移。
	//
	// 范围：[0,120]
	X float64 `json:"x"`
	// 画面在画布上的 y 坐标 (px)。
	//
	// 以画布左上角为原点，y 坐标为画面左上角相对于原点的纵向位移。
	//
	// 范围：[0,3840]
	Y float64 `json:"y"`
	// 画面的宽度 (px)。
	//
	// 范围：[120,3840]
	Width uint `json:"width"`
	// 画面的高度 (px)。
	//
	// 范围：[120,3840]
	Height uint `json:"height"`
	// 画面的图层编号。
	//  - 0 代表最下层的图层。
	//  - 100 代表最上层的图层。
	//
	// 范围：[0,100]
	ZOrder int `json:"zOrder"`
}

type CloudTranscoderCanvas struct {
	// 画面的宽度 (px)。
	//
	// 范围：[120,3840]
	Width uint `json:"width"`
	// 画面的高度 (px)。
	//
	// 范围：[120,3840]
	Height uint `json:"height"`
	// 画布的背景色。
	//
	// RGB 颜色值，以十进制数表示。
	//
	// 如 0 代表黑色，255 代表蓝色。
	//
	// 范围：[0,16777215]
	Color uint `json:"color"`
	// 画布背景图。
	//
	// 必须为合法 URL，且包含 jpg 或 png 后缀。
	//
	// 注意：如果不传值，则没有画布背景图。
	BackgroundImage string `json:"backgroundImage"`
	// 画布背景图的填充模式：
	//
	//  - "FILL"：在保持长宽比的前提下，缩放画面，并居中剪裁。
	//  - "FIT"：在保持长宽比的前提下，缩放画面，使其完整显示。
	//
	// 默认值："FILL"
	FillMode string `json:"fillMode"`
}

type CloudTranscoderWaterMark struct {
	// 水印图片的 URL。
	//
	// 必须为合法 URL，且包含 jpg 或 png 后缀。
	ImageURL string                 `json:"imageUrl"`
	Region   *CloudTranscoderRegion `json:"region"`
	// 画布背景图的填充模式：
	//
	//  - "FILL"：在保持长宽比的前提下，缩放画面，并居中剪裁。
	//  - "FIT"：在保持长宽比的前提下，缩放画面，使其完整显示。
	//
	// 默认值："FILL"
	FillMode string `json:"fillMode"`
}

type CloudTranscoderOutputAudioOption struct {
	// 转码输出的音频属性：
	//   - "AUDIO_PROFILE_DEFAULT"：48 kHz 采样率，音乐编码，单声道，编码码率最大值为 64 Kbps。
	//   - "AUDIO_PROFILE_SPEECH_STANDARD"：32 kHz 采样率，语音编码，单声道，编码码率最大值为 18 Kbps。
	//   - "AUDIO_PROFILE_MUSIC_STANDARD": 48 KHz 采样率，音乐编码，单声道，编码码率最大值为 64 Kbps。
	//   - "AUDIO_PROFILE_MUSIC_STANDARD_STEREO"：48 KHz 采样率，音乐编码，双声道，编码码率最大值为 80 Kbps。
	//   - "AUDIO_PROFILE_MUSIC_HIGH_QUALITY"：48 KHz 采样率，音乐编码，单声道，编码码率最大值为 96 Kbps。
	// 	 - "AUDIO_PROFILE_MUSIC_HIGH_QUALITY_STEREO"：48 KHz 采样率，音乐编码，双声道，编码码率最大值为 128 Kbps。
	//
	// 默认值："AUDIO_PROFILE_DEFAULT"
	ProfileType string `json:"profileType"`
}

type CloudTranscoderOutputVideoOption struct {
	// 	转码输出视频的帧率 (fps)。
	//
	// 范围：[1,30]
	//
	// 默认值：15
	FPS uint `json:"fps"`
	// 转码输出视频的 codec。取值包括：
	//  - "H264"：标准 H.264 编码。
	//  - "VP8"：标准 VP8 编码。
	Codec string `json:"codec"`
	// 	转码输出视频的码率。
	//
	// 范围：[1,10000]
	//
	// 注意：如果你不传值，声网会根据网络情况和其他视频属性自动设置视频码率。
	Bitrate uint `json:"bitrate"`
	// 画面的宽度 (px)。
	//
	// 范围：[120,3840]
	Width uint `json:"width"`
	// 画面的高度 (px)。
	//
	// 范围：[120,3840]
	Height uint `json:"height"`
	// 是否对转码输出的音视频流开启低码高清功能：
	//  - true：开启。
	//  - false：不开启。
	//
	// 默认值：false
	//
	// 注意：低码高清功能可以以较低的码率达到人眼感受的较高质量视频的效果。
	// 该功能会增加费用，开通前请联系声网销售咨询费用。
	LowBitrateHighQuality bool `json:"lowBitrateHighQuality"`
}

type CloudTranscoderOutput struct {
	Rtc         *CloudTranscoderRtc               `json:"rtc"`
	AudioOption *CloudTranscoderOutputAudioOption `json:"audioOption"`
	VideoOption *CloudTranscoderOutputVideoOption `json:"videoOption"`
}

type CreateSuccessResp struct {
	// 转码任务 ID，为 UUID，用于标识本次请求操作的 cloud transcoder
	TaskID string `json:"taskId"`
	// 转码任务创建时的 Unix 时间戳（秒）
	CreateTs int64 `json:"createTs"`
	// 转码任务的运行状态：
	//  - "IDLE": 任务未开始。
	//
	//  - "PREPARED": 任务已收到开启请求。
	//
	//  - "STARTING": 任务正在开启。
	//
	//  - "CREATED": 任务初始化完成。
	//
	//  - "STARTED": 任务已经启动。
	//
	//  - "IN_PROGRESS": 任务正在进行。
	//
	//  - "STOPPING": 任务正在停止。
	//
	//  - "STOPPED": 任务已经停止。
	//
	//  - "EXIT": 任务正常退出。
	//
	//  - "FAILURE_STOP": 任务异常退出。
	//
	// 注意：你可以用该字段监听任务的状态。
	Status string `json:"status"`
}

type CreateResp struct {
	Response
	SuccessResp CreateSuccessResp
}

func (c *Create) Do(ctx context.Context, tokenName string, payload *CreateReqBody) (*CreateResp, error) {
	path := c.buildPath(tokenName)
	responseData, err := c.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp CreateResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse CreateSuccessResp
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
