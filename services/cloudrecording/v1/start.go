package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO/agora-rest-client-go/core"
)

type Starter struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

const (
	IndividualMode = "individual"
	MixMode        = "mix"
	WebMode        = "web"
)

// BuildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/mode/{mode}/start
func (a *Starter) BuildPath(resourceID string, mode string) string {
	return a.prefixPath + "/resourceid/" + resourceID + "/mode/" + mode + "/start"
}

type StartReqBody struct {
	Cname         string              `json:"cname"`
	Uid           string              `json:"uid"`
	ClientRequest *StartClientRequest `json:"clientRequest"`
}

type StartClientRequest struct {
	Token                  string                  `json:"token,omitempty"`
	AppsCollection         *AppsCollection         `json:"appsCollection,omitempty"`
	RecordingConfig        *RecordingConfig        `json:"recordingConfig,omitempty"`
	TranscodeOptions       *TranscodeOptions       `json:"transcodeOptions,omitempty"`
	RecordingFileConfig    *RecordingFileConfig    `json:"recordingFileConfig,omitempty"`
	SnapshotConfig         *SnapshotConfig         `json:"snapshotConfig,omitempty"`
	StorageConfig          *StorageConfig          `json:"storageConfig,omitempty"`
	ExtensionServiceConfig *ExtensionServiceConfig `json:"extensionServiceConfig,omitempty"`
}

const (
	DefaultCombinationPolicy              = "default"
	PostPhoneTranscodingCombinationPolicy = "postphone_transcoding"
)

type AppsCollection struct {
	CombinationPolicy string `json:"combinationPolicy"`
}

type ChannelType int

const (
	CommunicationChannelType    ChannelType = 0
	LiveBroadcastingChannelType ChannelType = 1
)

type StreamType int

const (
	AudioStreamType      StreamType = 0
	VideoStreamType      StreamType = 1
	AudioVideoStreamType StreamType = 2
)

type StreamMode string

const (
	DefaultStreamMode  StreamMode = "default"
	StandardStreamMode StreamMode = "standard"
	OriginalStreamMode StreamMode = "original"
)

type DecryptionMode int

const (
	DefaultDecryptionMode   DecryptionMode = 0
	AES128XTSDecryptionMode DecryptionMode = 1
	AES128ECBDecryptionMode DecryptionMode = 2
	AES256XTSDecryptionMode DecryptionMode = 3
	SM4128ECBDecryptionMode DecryptionMode = 4
	AES128GCMDecryptionMode DecryptionMode = 5
	AES256GCMDecryptionMode DecryptionMode = 6
	AES128GCM2DecryptionMod DecryptionMode = 7
	AES256GCM2DecryptionMod DecryptionMode = 8
)

type RecordingConfig struct {
	ChannelType    int    `json:"channelType"`
	StreamTypes    int    `json:"streamTypes"`
	StreamMode     string `json:"streamMode,omitempty"`
	DecryptionMode int    `json:"decryptionMode,omitempty"`
	Secret         string `json:"secret,omitempty"`
	Salt           string `json:"salt,omitempty"`

	// AudioProfile 设置输出音频的采样率、码率、编码模式和声道数。目前仅适用于合流录制
	//
	// 0：(默认)48 kHz 采样率，音乐编码，单声道，编码码率约 48 Kbps
	//
	// 1：48 kHz 采样率，音乐编码，单声道，编码码率约 128 Kbps
	//
	// 2：48 kHz 采样率，音乐编码，双声道，编码码率约 192 Kbps
	AudioProfile int `json:"audioProfile,omitempty"`

	// VideoStreamType 设置订阅的视频流类型。如果频道中有用户开启了双流模式，
	// 你可以选择订阅视频大流或者小流。
	//
	// 0：(默认)视频大流，即高分辨率高码率的视频流
	//
	// 1：视频小流，即低分辨率低码率的视频流
	VideoStreamType int `json:"videoStreamType,omitempty"`

	// MaxIdleTime 最长空闲频道时间，单位为秒。默认值为 30。该值需大于等于 5，
	// 且小于等于 2,592,000，即 30 天。如果频道内无用户的状态持续超过该时间，
	// 录制程序会自动退出。退出后，再次调用 start 请求，会产生新的录制文件。
	//
	// * 通信场景下，如果频道内有用户，但用户没有发流，不算作无用户状态。
	//
	// * 直播场景下，如果频道内有观众但无主播，一旦无主播的状态超过 maxIdleTime，录制程序会自动退出。
	MaxIdleTime int `json:"maxIdleTime,omitempty"`

	TranscodingConfig *TranscodingConfig `json:"transcodingConfig,omitempty"`

	SubscribeAudioUIDs   []string `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUIDs []string `json:"unsubscribeAudioUids,omitempty"`
	SubscribeVideoUIDs   []string `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUIDs []string `json:"unsubscribeVideoUids,omitempty"`

	SubscribeUidGroup int `json:"subscribeUidGroup"`
}

type Container struct {
	Format string `json:"format"`
}

type TranscodeOptions struct {
	Container   *Container   `json:"container"`
	TransConfig *TransConfig `json:"transConfig"`
	Audio       *Audio       `json:"audio"`
}

type TransConfig struct {
	TransMode string `json:"transMode"`
}

type Audio struct {
	SampleRate string `json:"sampleRate"`
	BitRate    string `json:"bitrate"`
	Channels   string `json:"channels"`
}

type TranscodingConfig struct {
	Width                      int                `json:"width"`
	Height                     int                `json:"height"`
	FPS                        int                `json:"fps"`
	BitRate                    int                `json:"bitrate"`
	MaxResolutionUid           string             `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout           int                `json:"mixedVideoLayout"`
	BackgroundColor            string             `json:"backgroundColor,omitempty"`
	BackgroundImage            string             `json:"backgroundImage,omitempty"`
	DefaultUserBackgroundImage string             `json:"defaultUserBackgroundImage,omitempty"`
	LayoutConfig               []LayoutConfig     `json:"layoutConfig,omitempty"`
	BackgroundConfig           []BackgroundConfig `json:"backgroundConfig,omitempty"`
}

type LayoutConfig struct {
	UID        string  `json:"uid"`
	XAxis      float32 `json:"x_axis"`
	YAxis      float32 `json:"y_axis"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Alpha      float32 `json:"alpha"`
	RenderMode int     `json:"render_mode"`
}

type RecordingFileConfig struct {
	AvFileType []string `json:"avFileType"`
}

type SnapshotConfig struct {
	CaptureInterval int      `json:"captureInterval"`
	FileType        []string `json:"fileType"`
}
type StorageConfig struct {
	Vendor         int      `json:"vendor"`
	Region         int      `json:"region"`
	Bucket         string   `json:"bucket"`
	AccessKey      string   `json:"accessKey"`
	SecretKey      string   `json:"secretKey"`
	FileNamePrefix []string `json:"fileNamePrefix,omitempty"`
}

type ExtensionParams struct {
	SSE string `json:"sse"`
	Tag string `json:"tag"`
}

type ExtensionServiceConfig struct {
	ErrorHandlePolicy string             `json:"errorHandlePolicy"`
	APIVersion        string             `json:"apiVersion"`
	ExtensionServices []ExtensionService `json:"extensionServices"`
}

type ExtensionService struct {
	ServiceName       string        `json:"serviceName"`
	ErrorHandlePolicy string        `json:"errorHandlePolicy"`
	ServiceParam      *ServiceParam `json:"serviceParam"`
}

type Outputs struct {
	RtmpURL string `json:"rtmpUrl"`
}

type ServiceParam struct {
	Outputs          []Outputs `json:"outputs"`
	URL              string    `json:"url"`
	VideoBitRate     int       `json:"VideoBitrate"`
	VideoFPS         int       `json:"videoFps"`
	AudioProfile     int       `json:"audioProfile"`
	Mobile           bool      `json:"mobile"`
	VideoWidth       int       `json:"videoWidth"`
	VideoHeight      int       `json:"videoHeight"`
	MaxRecordingHour int       `json:"maxRecordingHour"`
	MaxVideoDuration int       `json:"MaxVideoDuration"`
	Onhold           bool      `json:"onhold"`
	ReadyTimeout     int       `json:"readyTimeout"`
}

type StarterResp struct {
	Response
	SuccessResp StartSuccessResp
}

type StartSuccessResp struct {
	CName      string `json:"cname"`
	UID        string `json:"uid"`
	ResourceId string `json:"resourceId"`
	SID        string `json:"sid"`
}

func (a *Starter) Do(ctx context.Context, resourceID string, mode string, payload *StartReqBody) (*StarterResp, error) {
	path := a.BuildPath(resourceID, mode)

	responseData, err := a.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp StarterResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResp StartSuccessResp
		if err = responseData.UnmarshallToTarget(&successResp); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResp
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
