package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Starter struct {
	module     string
	logger     core.Logger
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
func (s *Starter) BuildPath(resourceID string, mode string) string {
	return s.prefixPath + "/resourceid/" + resourceID + "/mode/" + mode + "/start"
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

type RecordingConfig struct {
	// ChannelType 频道场景。
	//
	// 目前支持以下几种频道场景：
	//
	// 0: 通信场景（默认）
	//
	// 1: 直播场景
	ChannelType int `json:"channelType"`

	// StreamTypes 订阅的媒体流类型
	//
	// 目前支持以下几种媒体流类型：
	//
	// 0:仅订阅音频。适用于智能语音审核场景
	//
	// 1:仅订阅视频
	//
	// 2:同时订阅音频和视频(默认)
	StreamTypes int `json:"streamTypes"`

	// StreamMode 媒体流的输出模式
	//
	// 目前支持以下几种输出模式：
	//
	// default:默认模式。录制过程中音频转码，分别生成 M3U8 音频索引文件和视频索引文件。
	//
	// standard:标准模式。声网推荐使用该模式。录制过程中音频转码，
	// 分别生成 M3U8 音频索引文件、视频索引文件和合并的音视频索引文件。如果在 Web 端使用 VP8 编码，则生成一个合并的 MPD 音视频索引文件。
	//
	// original:原始编码模式。
	//适用于单流音频不转码录制。仅订阅音频时（streamTypes 为 0）时该字段生效，录制过程中音频不转码，生成 M3U8 音频索引文件。
	StreamMode string `json:"streamMode,omitempty"`

	// DecryptionMode 音频流的解密模式
	//
	// 下面是支持的解密模式：
	//0：不加密。
	//
	//1：AES_128_XTS 加密模式。128 位 AES 加密，XTS 模式。
	//
	//2：AES_128_ECB 加密模式。128 位 AES 加密，ECB 模式。
	//
	//3：AES_256_XTS 加密模式。256 位 AES 加密，XTS 模式。
	//
	//4：SM4_128_ECB 加密模式。128 位 SM4 加密，ECB 模式。
	//
	//5：AES_128_GCM 加密模式。128 位 AES 加密，GCM 模式。
	//
	//6：AES_256_GCM 加密模式。256 位 AES 加密，GCM 模式。
	//
	//7：AES_128_GCM2 加密模式。128 位 AES 加密，GCM 模式。相比于 AES_128_GCM 加密模式，AES_128_GCM2 加密模式安全性更高且需要设置密钥和盐。
	//
	//8：AES_256_GCM2 加密模式。256 位 AES 加密，GCM 模式。相比于 AES_256_GCM 加密模式，AES_256_GCM2 加密模式安全性更高且需要设置密钥和盐。
	DecryptionMode int `json:"decryptionMode,omitempty"`

	// Secret 与加解密相关的密钥
	//
	// 仅需在 decryptionMode 非 0 时设置。
	Secret string `json:"secret,omitempty"`

	// Salt 与加解密相关的盐
	//
	// Base64 编码、32 位字节。仅需在 decryptionMode 为 7 或 8 时设置。
	Salt string `json:"salt,omitempty"`

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
	// 通信场景下，如果频道内有用户，但用户没有发流，不算作无用户状态。
	//
	// 直播场景下，如果频道内有观众但无主播，一旦无主播的状态超过 maxIdleTime，录制程序会自动退出。
	MaxIdleTime int `json:"maxIdleTime,omitempty"`

	// TranscodingConfig 转码输出的视频配置项
	//
	// 配置参考 https://doc.shengwang.cn/doc/cloud-recording/restful/user-guide/mix-mode/set-output-video-profile
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
	// Width 视频的宽度，单位为像素。
	//
	// width 和 height 的乘积不能超过 1920 × 1080。
	//
	// 默认值：360
	Width int `json:"width,omitempty"`

	// Height 视频的高度，单位为像素。
	//
	// width 和 height 的乘积不能超过 1920 × 1080。
	//
	// 默认值：640
	Height int `json:"height,omitempty"`

	// Fps 视频的帧率,单位 fps。
	//
	// 默认值：15
	FPS int `json:"fps,omitempty"`

	// BitRate 视频的码率，单位为 Kbps。
	//
	// 默认值：500
	BitRate int `json:"bitrate,omitempty"`

	// MaxResolutionUid 仅需在垂直布局下设置。
	//
	//指定显示大视窗画面的用户 UID。字符串内容的整型取值范围 1 到 (232-1)，且不可设置为 0。
	MaxResolutionUid string `json:"maxResolutionUid,omitempty"`

	// MixedVideoLayout 视频合流布局
	//
	//0：悬浮布局。第一个加入频道的用户在屏幕上会显示为大视窗，铺满整个画布，其他用户的视频画面会显示为小视窗，从下到上水平排列，最多 4 行，每行 4 个画面，最多支持共 17 个画面。
	//
	//1：自适应布局。根据用户的数量自动调整每个画面的大小，每个用户的画面大小一致，最多支持 17 个画面。
	//
	//2：垂直布局。指定 maxResolutionUid 在屏幕左侧显示大视窗画面，其他用户的小视窗画面在右侧垂直排列，最多两列，一列 8 个画面，最多支持共 17 个画面。
	//
	//3：自定义布局。由你在 layoutConfig 字段中自定义合流布局。
	//
	// 默认值：0
	MixedVideoLayout int `json:"mixedVideoLayout,omitempty"`

	// BackgroundColor 视频画布的背景颜色。
	//
	//支持 RGB 颜色表，字符串格式为 # 号和 6 个十六进制数。
	//
	//默认值 "#000000"，代表黑色。
	BackgroundColor string `json:"backgroundColor,omitempty"`

	// BackgroundImage 视频画布的背景图的 URL,背景图的显示模式为裁剪模式。
	//
	//裁剪模式：优先保证画面被填满。背景图尺寸等比缩放，直至整个画面被背景图填满。如果背景图长宽与显示窗口不同，则背景图会按照画面设置的比例进行周边裁剪后填满画面。
	BackgroundImage string `json:"backgroundImage,omitempty"`

	// DefaultUserBackgroundImage 默认的用户画面背景图的 URL
	DefaultUserBackgroundImage string `json:"defaultUserBackgroundImage,omitempty"`

	// LayoutConfig 用户的合流画面布局。
	//
	//由每个用户对应的布局画面设置组成的数组，支持最多 17 个用户。
	LayoutConfig []LayoutConfig `json:"layoutConfig,omitempty"`

	// BackgroundConfig 用户的背景图设置
	BackgroundConfig []BackgroundConfig `json:"backgroundConfig,omitempty"`
}

type LayoutConfig struct {
	// UID 字符串内容为待显示在该区域的用户的 UID，32 位无符号整数。
	//
	//如果不指定 UID，会按照用户加入频道的顺序自动匹配 layoutConfig 中的画面设置。
	UID string `json:"uid"`

	// XAxis 屏幕里该画面左上角的横坐标的相对值，精确到小数点后六位。
	//
	//从左到右布局，0.0 在最左端，1.0 在最右端。该字段也可以设置为整数 0 或 1。
	//
	// 0<=XAxis<=1
	XAxis float32 `json:"x_axis"`

	// YAxis 屏幕里该画面左上角的纵坐标的相对值，精确到小数点后六位。
	//
	//屏幕里该画面左上角的纵坐标的相对值，精确到小数点后六位。
	//
	//从上到下布局，0.0 在最上端，1.0 在最下端。该字段也可以设置为整数 0 或 1。
	//
	// 0<=YAxis<=1
	YAxis float32 `json:"y_axis"`

	// Width 该画面宽度的相对值，精确到小数点后六位。该字段也可以设置为整数 0 或 1。
	//
	// 0<=Width<=1
	Width float32 `json:"width"`

	// Height 该画面高度的相对值，精确到小数点后六位。该字段也可以设置为整数 0 或 1。
	//
	// 0<=Height<=1
	Height float32 `json:"height"`

	// Alpha 图像的透明度。精确到小数点后六位。0.0 表示图像为透明的，1.0 表示图像为完全不透明的。
	//
	// 0<=Alpha<=1
	//
	// 默认值：1
	Alpha float32 `json:"alpha,omitempty"`

	// RenderMode 画面的渲染模式。
	//
	//0：裁剪模式。优先保证画面被填满。视频尺寸等比缩放，直至整个画面被视频填满。如果视频长宽与显示窗口不同，则视频流会按照画面设置的比例进行周边裁剪后填满画面。
	//
	//1：缩放模式。优先保证视频内容全部显示。视频尺寸等比缩放，直至视频窗口的一边与画面边框对齐。如果视频尺寸与画面尺寸不一致，在保持长宽比的前提下，将视频进行缩放后填满画面，缩放后的视频四周会有一圈黑边。
	//
	// 默认值：0
	RenderMode int `json:"render_mode"`
}

type RecordingFileConfig struct {
	// AvFileType 录制生成视频的文件类型
	//
	//"hls"：默认值。M3U8 和 TS 文件。
	//
	//"mp4"：MP4 文件。
	//
	// 注意：
	//
	//单流录制模式下，且非仅截图情况，使用默认值即可。
	//
	//合流录制和页面录制模式下，你需设为 ["hls","mp4"]或者["hls"]。仅设为 ["mp4"] 会收到报错。设置后，录制文件行为如下：
	//
	//合流录制模式：录制服务会在当前 MP4 文件时长超过约 2 小时或文件大小超过约 2 GB 左右时，创建一个新的 MP4 文件。
	//
	//页面录制模式：录制服务会在当前 MP4 文件时长超过 maxVideoDuration 时，创建一个新的 MP4 文件。
	AvFileType []string `json:"avFileType"`
}

type SnapshotConfig struct {
	// CaptureInterval 云端录制定期截图的截图周期。单位为秒。
	//
	// 5<=CaptureInterval<=300
	// 默认值：10
	CaptureInterval int `json:"captureInterval,omitempty"`

	// FileType 截图的文件格式
	//
	// 目前只支持 ["jpg"]，即生成 JPG 格式的截图文件。
	FileType []string `json:"fileType"`
}
type StorageConfig struct {
	// Vendor 第三方云存储平台。目前支持的云存储服务商有：
	//
	// 1:AWS S3
	//
	// 2:阿里云
	//
	// 3:腾讯云
	//
	// 5:Microsoft Azure
	//
	// 6:谷歌云
	//
	// 7:华为云
	//
	// 8:百度智能云
	Vendor int `json:"vendor"`

	// Region 第三方云存储指定的地区信息，详情见 https://doc.shengwang.cn/api-ref/cloud-recording/restful/region-vendor
	//
	// 注意：为确保录制文件上传的成功率和实时性，第三方云存储的 region 与你发起请求的应用服务器必须在同一个区域中。例如：你发起请求的 App 服务器在中国大陆地区，则第三方云存储需要设置为中国大陆区域内。
	Region int `json:"region"`

	// Bucket 第三方云存储的 bucket
	//
	// Bucket 名称需要符合对应第三方云存储服务的命名规则。
	Bucket string `json:"bucket"`

	// AccessKey 第三方云存储的 Access Key
	//
	// 如需延时转码，则访问密钥必须具备读写权限；否则建议只需提供写权限。
	AccessKey string `json:"accessKey"`

	// SecretKey 第三方云存储的 Secret Key
	SecretKey string `json:"secretKey"`

	// FileNamePrefix 录制文件的文件名前缀
	//
	// 录制文件在第三方云存储中的存储位置，与录制文件名前缀有关。
	// 如果设为 ["directory1","directory2"]，那么录制文件名前缀为 "directory1/directory2/"，
	// 即录制文件名为 directory1/directory2/xxx.m3u8。前缀长度（包括斜杠）不得超过 128 个字符。字符串中不得出现斜杠、下划线、括号等符号字符。
	//
	// 以下为支持的字符集范围：
	//
	// 26 个小写英文字母 a~z
	//
	// 26 个大写英文字母 A~Z
	//
	// 10 个数字 0-9
	FileNamePrefix []string `json:"fileNamePrefix,omitempty"`

	// ExtensionParams 第三方云存储服务会按照该字段设置对已上传的录制文件进行加密和打标签
	ExtensionParams *ExtensionParams `json:"extensionParams,omitempty"`
}

type ExtensionParams struct {
	// SSE 加密模式
	// 设置该字段后，第三方云存储服务会按照该加密模式将已上传的录制文件进行加密。该字段仅适用于 Amazon S3
	//
	// kms:KMS 加密。
	//
	// aes256:AES256 加密
	SSE string `json:"sse"`

	// Tag 标签内容
	// 设置该字段后，第三方云存储服务会按照该标签内容将已上传的录制文件进行打标签操作。该字段仅适用于阿里云和 Amazon S3
	Tag string `json:"tag"`
}

type ExtensionServiceConfig struct {
	ErrorHandlePolicy string             `json:"errorHandlePolicy,omitempty"`
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
	Outputs          []Outputs `json:"outputs,omitempty"`
	URL              string    `json:"url"`
	VideoBitRate     int       `json:"VideoBitrate,omitempty"`
	VideoFPS         int       `json:"videoFps,omitempty"`
	AudioProfile     int       `json:"audioProfile"`
	Mobile           bool      `json:"mobile,omitempty"`
	VideoWidth       int       `json:"videoWidth"`
	VideoHeight      int       `json:"videoHeight"`
	MaxRecordingHour int       `json:"maxRecordingHour"`
	MaxVideoDuration int       `json:"maxVideoDuration,omitempty"`
	Onhold           bool      `json:"onhold,omitempty"`
	ReadyTimeout     int       `json:"readyTimeout,omitempty"`
}

type StarterResp struct {
	Response
	SuccessResp StartSuccessResp
}

type StartSuccessResp struct {
	Cname      string `json:"cname"`
	UID        string `json:"uid"`
	ResourceId string `json:"resourceId"`
	SID        string `json:"sid"`
}

func (s *Starter) Do(ctx context.Context, resourceID string, mode string, payload *StartReqBody) (*StarterResp, error) {
	path := s.BuildPath(resourceID, mode)

	responseData, err := s.doRESTWithRetry(ctx, path, http.MethodPost, payload)
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

const retryCount = 3

func (s *Starter) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*core.BaseResponse, error) {
	var (
		resp  *core.BaseResponse
		err   error
		retry int
	)

	err = core.RetryDo(func(retryCount int) error {
		var doErr error

		resp, doErr = s.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return core.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 410:
			s.logger.Debugf(ctx, s.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return core.NewRetryErr(
				false,
				core.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			)
		default:
			s.logger.Debugf(ctx, s.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retry >= retryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		s.logger.Debugf(ctx, s.module, "http request err:%s", err)
		retry++
	})

	return resp, err
}
