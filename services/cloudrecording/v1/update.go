package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Update struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/update
func (u *Update) buildPath(resourceID string, sid string, mode string) string {
	return u.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/update"
}

type UpdateReqBody struct {
	Cname         string               `json:"cname"`
	Uid           string               `json:"uid"`
	ClientRequest *UpdateClientRequest `json:"clientRequest"`
}

type UpdateClientRequest struct {
	StreamSubscribe    *UpdateStreamSubscribe    `json:"streamSubscribe,omitempty"`
	WebRecordingConfig *UpdateWebRecordingConfig `json:"webRecordingConfig,omitempty"`
	RtmpPublishConfig  *UpdateRtmpPublishConfig  `json:"rtmpPublishConfig,omitempty"`
}

// UpdateStreamSubscribe 更新订阅名单
type UpdateStreamSubscribe struct {
	// AudioUidList 音频订阅名单
	AudioUidList *UpdateAudioUIDList `json:"audioUidList,omitempty"`

	// VideoUidList 视频订阅名单
	VideoUidList *UpdateVideoUIDList `json:"videoUidList,omitempty"`
}

// UpdateAudioUIDList 音频订阅名单
type UpdateAudioUIDList struct {
	// SubscribeAudioUIDs 指定订阅哪几个 UID 的音频流
	//
	// 如需订阅全部 UID 的音频流，则无需设置该字段。数组长度不得超过 32，不推荐使用空数组。该字段和 unsubscribeAudioUids 只能设一个。
	//
	// 注意：
	//
	// 该字段仅适用于 streamTypes 设为音频，或音频和视频的情况。
	//
	// 如果你设置了音频的订阅名单，但没有设置视频的订阅名单，云端录制服务不会订阅任何视频流。反之亦然。
	//
	// 设为 ["#allstream#"] 可订阅频道内所有 UID 的音频流。
	SubscribeAudioUIDs []string `json:"subscribeAudioUids,omitempty"`

	// UnsubscribeAudioUIDs 指定取消订阅哪几个 UID 的音频流
	//
	// 云端录制会订阅频道内除指定 UID 外所有 UID 的音频流。数组长度不得超过 32，不推荐使用空数组。该字段和 subscribeAudioUids 只能设一个。
	//
	UnsubscribeAudioUIDs []string `json:"unSubscribeAudioUids,omitempty"`
}

// UpdateVideoUIDList 视频订阅名单
type UpdateVideoUIDList struct {
	// SubscribeVideoUIDs 指定订阅哪几个 UID 的视频流
	//
	// 如需订阅全部 UID 的视频流，则无需设置该字段。数组长度不得超过 32，不推荐使用空数组。该字段和 unsubscribeVideoUids 只能设一个。
	//
	// 注意：
	//
	// 该字段仅适用于 streamTypes 设为视频，或音频和视频的情况。
	//
	// 如果你设置了视频的订阅名单，但没有设置音频的订阅名单，云端录制服务不会订阅任何音频流。反之亦然。
	//
	// 设为 ["#allstream#"] 可订阅频道内所有 UID 的视频流。
	SubscribeVideoUIDs []string `json:"subscribeVideoUids,omitempty"`

	// UnsubscribeVideoUIDs 指定取消订阅哪几个 UID 的视频流
	//
	// 云端录制会订阅频道内除指定 UID 外所有 UID 的视频流。数组长度不得超过 32，不推荐使用空数组。该字段和 subscribeVideoUids 只能设一个。
	UnsubscribeVideoUIDs []string `json:"unSubscribeVideoUids,omitempty"`
}

// UpdateWebRecordingConfig 用于更新页面录制配置项。
type UpdateWebRecordingConfig struct {
	// Onhold 是否在启动页面录制任务时暂停页面录制。
	//
	// true：在启动页面录制任务时暂停页面录制。开启页面录制任务后立即暂停录制，录制服务会打开并渲染待录制页面，但不生成切片文件。
	//
	// false：启动页面录制任务并进行页面录制。
	//
	// 建议你按照如下流程使用 onhold 字段：
	//
	// 调用 start 方法时将 onhold 设为 true，开启并暂停页面录制，自行判断页面录制开始的合适时机。
	//
	// 调用 update 并将 onhold 设为 false，继续进行页面录制。如果需要连续调用 update 方法暂停或继续页面录制，请在收到上一次 update 响应后再进行调用，否则可能导致请求结果与预期不一致。
	//
	// 默认值：false
	Onhold bool `json:"onhold"`
}

// UpdateRtmpPublishConfig 用于更新转推页面录制到 CDN 的配置项。
type UpdateRtmpPublishConfig struct {
	Outputs []UpdateOutput `json:"outputs"`
}

type UpdateOutput struct {
	// RtmpURL CDN 推流 URL。
	//
	// 注意：
	//
	// URL 仅支持 RTMP 和 RTMPS 协议。
	//
	// 支持的最大转推 CDN 路数为 1。
	RtmpURL string `json:"rtmpUrl"`
}

type UpdateResp struct {
	Response
	SuccessResponse UpdateSuccessResp
}

type UpdateSuccessResp struct {
	ResourceId string `json:"resourceId"`
	Sid        string `json:"sid"`
	UID        string `json:"uid"`
	Cname      string `json:"cname"`
}

func (u *Update) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateReqBody) (*UpdateResp, error) {
	path := u.buildPath(resourceID, sid, mode)

	responseData, err := u.client.DoREST(ctx, path, http.MethodPost, payload)
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
		resp.SuccessResponse = successResponse
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
