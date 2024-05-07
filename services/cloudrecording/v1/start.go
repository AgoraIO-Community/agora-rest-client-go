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

type Start struct {
	forwardedRegionPrefix core.ForwardedReginPrefix
	module                string
	logger                core.Logger
	client                core.Client
	prefixPath            string // /v1/apps/{appid}/cloud_recording
}

const (
	// IndividualMode 云端录制模式：单流录制
	IndividualMode = "individual"

	// MixMode 云端录制模式：合流录制
	MixMode = "mix"

	// WebMode 云端录制模式：页面录制
	WebMode = "web"
)

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/mode/{mode}/start
func (s *Start) buildPath(resourceID string, mode string) string {
	return string(s.forwardedRegionPrefix) + s.prefixPath + "/resourceid/" + resourceID + "/mode/" + mode + "/start"
}

type StartReqBody struct {
	// Cname 录制的频道名
	Cname string `json:"cname"`

	// Uid 字符串内容为云端录制服务在 RTC 频道内使用的 UID，用于标识频道内的录制服务。
	Uid string `json:"uid"`

	// ClientRequest 客户端请求
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
	PostPhoneTranscodingCombinationPolicy = "postpone_transcoding"
)

// AppsCollection 应用配置项
type AppsCollection struct {
	// CombinationPolicy 各云端录制应用的组合方式
	//
	// postpone_transcoding：如需延时转码或延时混音，则选用此种方式。
	//
	// default：除延时转码和延时混音外，均选用此种方式。
	//
	// 默认值：default
	CombinationPolicy string `json:"combinationPolicy,omitempty"`
}

type RecordingConfig struct {
	// ChannelType 频道场景。
	//
	// 目前支持以下几种频道场景：
	//
	// 0: 通信场景
	//
	// 1: 直播场景
	//
	// 默认值：0
	ChannelType int `json:"channelType"`

	// StreamTypes 订阅的媒体流类型
	//
	// 目前支持以下几种媒体流类型：
	//
	// 0:仅订阅音频。适用于智能语音审核场景
	//
	// 1:仅订阅视频
	//
	// 2:同时订阅音频和视频
	//
	// 默认值：2
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
	//
	// 默认值：default
	StreamMode string `json:"streamMode,omitempty"`

	// DecryptionMode 音频流的解密模式
	//
	// 下面是支持的解密模式：
	// 0：不加密。
	//
	// 1：AES_128_XTS 加密模式。128 位 AES 加密，XTS 模式。
	//
	// 2：AES_128_ECB 加密模式。128 位 AES 加密，ECB 模式。
	//
	// 3：AES_256_XTS 加密模式。256 位 AES 加密，XTS 模式。
	//
	// 4：SM4_128_ECB 加密模式。128 位 SM4 加密，ECB 模式。
	//
	// 5：AES_128_GCM 加密模式。128 位 AES 加密，GCM 模式。
	//
	// 6：AES_256_GCM 加密模式。256 位 AES 加密，GCM 模式。
	//
	// 7：AES_128_GCM2 加密模式。128 位 AES 加密，GCM 模式。相比于 AES_128_GCM 加密模式，AES_128_GCM2 加密模式安全性更高且需要设置密钥和盐。
	//
	// 8：AES_256_GCM2 加密模式。256 位 AES 加密，GCM 模式。相比于 AES_256_GCM 加密模式，AES_256_GCM2 加密模式安全性更高且需要设置密钥和盐。
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

	// MaxIdleTime 最大频道空闲时间。
	// 单位为秒。最大值不超过 30 天。超出最大频道空闲时间后，录制服务会自动退出。录制服务退出后，如果你再次发起 start 请求，会产生新的录制文件。
	//
	// 频道空闲：直播频道内无任何主播，或通信频道内无任何用户。
	//
	// 5 <= MaxIdleTime <= 2592000
	//
	// 默认值：30
	MaxIdleTime int `json:"maxIdleTime,omitempty"`

	// TranscodingConfig 转码输出的视频配置项
	//
	// 配置参考 https://doc.shengwang.cn/doc/cloud-recording/restful/user-guide/mix-mode/set-output-video-profile
	TranscodingConfig *TranscodingConfig `json:"transcodingConfig,omitempty"`

	// SubscribeUIDs 指定订阅哪几个 UID 的音频流。
	//
	// 如需订阅全部 UID 的音频流，则无需设置该字段。数组长度不得超过 32，不推荐使用空数组。该字段和 UnsubscribeAudioUIDs 只能设一个。
	//
	// 注意：
	//
	// 1.该字段仅适用于 streamTypes 设为音频，或音频和视频的情况。
	//
	// 2.如果你设置了音频的订阅名单，但没有设置视频的订阅名单，云端录制服务不会订阅任何视频流。反之亦然。
	//
	// 3.设为 ["#allstream#"] 可订阅频道内所有 UID 的音频流。
	SubscribeAudioUIDs []string `json:"subscribeAudioUids,omitempty"`

	// UnsubscribeAudioUIDs 指定不订阅哪几个 UID 的音频流。
	//
	// 云端录制会订阅频道内除指定 UID 外所有 UID 的音频流。数组长度不得超过 32，不推荐使用空数组。该字段和 SubscribeAudioUIDs 只能设一个。
	UnsubscribeAudioUIDs []string `json:"unsubscribeAudioUids,omitempty"`

	// SubscribeVideoUIDs 指定订阅哪几个 UID 的视频流。
	//
	// 如需订阅全部 UID 的视频流，则无需设置该字段。数组长度不得超过 32，不推荐使用空数组。该字段和 UnsubscribeVideoUIDs 只能设一个。
	//
	// 注意：
	//
	// 1.该字段仅适用于 streamTypes 设为视频，或音频和视频的情况。
	//
	// 2.如果你设置了视频的订阅名单，但没有设置音频的订阅名单，云端录制服务不会订阅任何音频流。反之亦然。
	//
	// 3.设为 ["#allstream#"] 可订阅频道内所有 UID 的视频流。
	SubscribeVideoUIDs []string `json:"subscribeVideoUids,omitempty"`

	// UnsubscribeVideoUIDs 指定不订阅哪几个 UID 的视频流。
	//
	// 云端录制会订阅频道内除指定 UID 外所有 UID 的视频流。数组长度不得超过 32，不推荐使用空数组。该字段和 subscribeVideoUids 只能设一个。
	UnsubscribeVideoUIDs []string `json:"unsubscribeVideoUids,omitempty"`

	// SubscribeUidGroup 预估的订阅人数峰值。
	//
	// 枚举值：
	//
	// 0：1 到 2 个 UID。
	//
	// 1：3 到 7 个 UID。
	//
	// 2：8 到 12 个 UID。
	//
	// 3：13 到 17 个 UID。
	//
	// 4：18 到 32 个 UID。
	//
	// 5：33 到 49 个 UID。
	SubscribeUidGroup int `json:"subscribeUidGroup"`
}

type Container struct {
	// Format 文件的容器格式，支持如下取值：
	//
	// "mp4"：延时转码时的默认格式。MP4 格式。
	//
	// "mp3"：延时混音时的默认格式。MP3 格式。
	//
	// "m4a"：M4A 格式。
	//
	// "aac"：AAC 格式。
	//
	// 注意：延时转码暂时只能设为 MP4 格式。
	Format string `json:"format"`
}

// TranscodeOptions 延时转码或延时混音下，生成的录制文件的配置项。
type TranscodeOptions struct {
	Container   *Container   `json:"container"`
	TransConfig *TransConfig `json:"transConfig"`
	Audio       *Audio       `json:"audio"`
}

type TransConfig struct {
	// TransMode 模式
	//
	// "postponeTranscoding"：延时转码。
	//
	// "audioMix"：延时混音。
	TransMode string `json:"transMode"`
}

// Audio 文件的音频属性
type Audio struct {
	// SampleRate 音频的采样率，单位为 Hz，支持如下取值：
	//
	// "48000"：48 kHz。
	//
	// "32000"：32 kHz。
	//
	// "16000"：16 kHz。
	SampleRate string `json:"sampleRate"`

	// BitRate 音频的码率，单位为 Kbps
	//
	// 默认值：4800
	BitRate string `json:"bitrate"`

	// Channels 音频声道数，支持如下取值：
	//
	// "1"：单声道。
	//
	// "2"：双声道。
	//
	// 默认值：2
	Channels string `json:"channels"`
}

// TranscodingConfig 转码输出的视频配置项
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
	// 指定显示大视窗画面的用户 UID。字符串内容的整型取值范围 1 到 (232-1)，且不可设置为 0。
	MaxResolutionUid string `json:"maxResolutionUid,omitempty"`

	// MixedVideoLayout 视频合流布局
	//
	// 0：悬浮布局。第一个加入频道的用户在屏幕上会显示为大视窗，铺满整个画布，其他用户的视频画面会显示为小视窗，从下到上水平排列，最多 4 行，每行 4 个画面，最多支持共 17 个画面。
	//
	// 1：自适应布局。根据用户的数量自动调整每个画面的大小，每个用户的画面大小一致，最多支持 17 个画面。
	//
	// 2：垂直布局。指定 maxResolutionUid 在屏幕左侧显示大视窗画面，其他用户的小视窗画面在右侧垂直排列，最多两列，一列 8 个画面，最多支持共 17 个画面。
	//
	// 3：自定义布局。由你在 layoutConfig 字段中自定义合流布局。
	//
	// 默认值：0
	MixedVideoLayout int `json:"mixedVideoLayout,omitempty"`

	// BackgroundColor 视频画布的背景颜色。
	//
	// 支持 RGB 颜色表，字符串格式为 # 号和 6 个十六进制数。
	//
	// 默认值 "#000000"，代表黑色。
	BackgroundColor string `json:"backgroundColor,omitempty"`

	// BackgroundImage 视频画布的背景图的 URL,背景图的显示模式为裁剪模式。
	//
	// 裁剪模式：优先保证画面被填满。背景图尺寸等比缩放，直至整个画面被背景图填满。如果背景图长宽与显示窗口不同，则背景图会按照画面设置的比例进行周边裁剪后填满画面。
	BackgroundImage string `json:"backgroundImage,omitempty"`

	// DefaultUserBackgroundImage 默认的用户画面背景图的 URL
	DefaultUserBackgroundImage string `json:"defaultUserBackgroundImage,omitempty"`

	// LayoutConfig 用户的合流画面布局。
	//
	// 由每个用户对应的布局画面设置组成的数组，支持最多 17 个用户。
	LayoutConfig []LayoutConfig `json:"layoutConfig,omitempty"`

	// BackgroundConfig 用户的背景图设置
	BackgroundConfig []BackgroundConfig `json:"backgroundConfig,omitempty"`
}

type LayoutConfig struct {
	// UID 字符串内容为待显示在该区域的用户的 UID，32 位无符号整数。
	//
	// 如果不指定 UID，会按照用户加入频道的顺序自动匹配 layoutConfig 中的画面设置。
	UID string `json:"uid"`

	// XAxis 屏幕里该画面左上角的横坐标的相对值，精确到小数点后六位。
	//
	// 从左到右布局，0.0 在最左端，1.0 在最右端。该字段也可以设置为整数 0 或 1。
	//
	// 0<=XAxis<=1
	XAxis float32 `json:"x_axis"`

	// YAxis 屏幕里该画面左上角的纵坐标的相对值，精确到小数点后六位。
	//
	// 屏幕里该画面左上角的纵坐标的相对值，精确到小数点后六位。
	//
	// 从上到下布局，0.0 在最上端，1.0 在最下端。该字段也可以设置为整数 0 或 1。
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
	// 0：裁剪模式。优先保证画面被填满。视频尺寸等比缩放，直至整个画面被视频填满。如果视频长宽与显示窗口不同，则视频流会按照画面设置的比例进行周边裁剪后填满画面。
	//
	// 1：缩放模式。优先保证视频内容全部显示。视频尺寸等比缩放，直至视频窗口的一边与画面边框对齐。如果视频尺寸与画面尺寸不一致，在保持长宽比的前提下，将视频进行缩放后填满画面，缩放后的视频四周会有一圈黑边。
	//
	// 默认值：0
	RenderMode int `json:"render_mode"`
}

type RecordingFileConfig struct {
	// AvFileType 录制生成视频的文件类型
	//
	// "hls"：默认值。M3U8 和 TS 文件。
	//
	// "mp4"：MP4 文件。
	//
	// 注意：
	//
	// 单流录制模式下，且非仅截图情况，使用默认值即可。
	//
	// 合流录制和页面录制模式下，你需设为 ["hls","mp4"]或者["hls"]。仅设为 ["mp4"] 会收到报错。设置后，录制文件行为如下：
	//
	// 合流录制模式：录制服务会在当前 MP4 文件时长超过约 2 小时或文件大小超过约 2 GB 左右时，创建一个新的 MP4 文件。
	//
	// 页面录制模式：录制服务会在当前 MP4 文件时长超过 maxVideoDuration 时，创建一个新的 MP4 文件。
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

// ExtensionServiceConfig 扩展服务配置项
type ExtensionServiceConfig struct {
	// ErrorHandlePolicy 错误处理策略。
	//
	// 默认且仅可设为 "error_abort"，表示当扩展服务发生错误后，订阅和云端录制的其他非扩展服务都停止。
	//
	// 默认值：error_abort
	ErrorHandlePolicy string `json:"errorHandlePolicy,omitempty"`

	ExtensionServices []ExtensionService `json:"extensionServices"`
}

// ExtensionService 扩展服务
type ExtensionService struct {
	// ServiceName 扩展服务的名称
	//
	// web_recorder_service：代表扩展服务为页面录制。
	//
	// rtmp_publish_service：代表扩展服务为转推页面录制到 CDN。
	ServiceName string `json:"serviceName"`

	// ErrorHandlePolicy 扩展服务内的错误处理策略
	//
	// "error_abort"：页面录制时默认且只能为该值。表示当前扩展服务出错时，停止其他扩展服务。
	//
	// "error_ignore"：转推页面录制到 CDN 时默认且只能为该值。表示当前扩展服务出错时，其他扩展服务不受影响。
	//
	// 如果页面录制服务或录制上传服务异常，那么推流到 CDN 失败，因此页面录制服务出错会影响转推页面录制到 CDN 服务。
	//
	// 转推到 CDN 的过程发生异常时，页面录制不受影响。
	ErrorHandlePolicy string `json:"errorHandlePolicy"`

	// ServiceParam 扩展服务的参数
	ServiceParam *ServiceParam `json:"serviceParam"`
}

type Outputs struct {
	// RtmpURL CDN 推流地址
	RtmpURL string `json:"rtmpUrl"`
}

type ServiceParam struct {
	// Outputs 转推页面录制到 CDN 时需设置如下字段
	Outputs []Outputs `json:"outputs,omitempty"`

	// URL 待录制页面的地址
	URL string `json:"url"`

	// VideoBitRate 输出视频的码率，单位为 Kbps
	//
	// 针对不同的输出视频分辨率，videoBitrate 的默认值不同：
	//
	// 输出视频分辨率大于或等于 1280 × 720：默认值为 2000。
	//
	// 输出视频分辨率小于 1280 × 720：默认值为 1500。
	VideoBitRate int `json:"VideoBitrate,omitempty"`

	// VideoFPS 输出视频的帧率，单位为 fps
	//
	//  5 <= VideoFPS <=60
	//
	// 默认值：15
	VideoFPS int `json:"videoFps,omitempty"`

	// AudioProfile 输出音频的采样率、码率、编码模式和声道数。
	//
	// 0：48 kHz 采样率，音乐编码，单声道，编码码率约 48 Kbps。
	//
	// 1：48 kHz 采样率，音乐编码，单声道，编码码率约 128 Kbps。
	//
	// 2：48 kHz 采样率，音乐编码，双声道，编码码率约 192 Kbps。
	AudioProfile int `json:"audioProfile"`

	// Mobile 是否开启移动端网页模式
	//
	// true：开启。开启后，录制服务使用移动端网页渲染模式录制当前页面。
	//
	// false：（默认）不开启。
	Mobile bool `json:"mobile,omitempty"`

	// VideoWidth 视频的宽度，单位为像素。
	//
	// videoWidth 和 videoHeight 的乘积需小于等于 1920 × 1080
	VideoWidth int `json:"videoWidth"`

	// VideoHeight 视频的高度，单位为像素。
	//
	// videoWidth 和 videoHeight 的乘积需小于等于 1920 × 1080
	VideoHeight int `json:"videoHeight"`

	// MaxRecordingHour 页面录制的最大时长，单位为小时
	//
	// 超出该值后页面录制会自动停止
	//
	// 计费相关：页面录制停止前会持续计费，因此请根据实际业务情况设置合理的值或主动停止页面录制。
	//
	// 1 <= MaxRecordingHour <= 720
	MaxRecordingHour int `json:"maxRecordingHour"`

	// MaxVideoDuration 页面录制生成的 MP4 切片文件的最大时长，单位为分钟
	//
	// 页面录制过程中，录制服务会在当前 MP4 文件时长超过约 maxVideoDuration 左右时创建一个新的 MP4 切片文件
	//
	// 30 <= MaxVideoDuration <= 240
	//
	// 默认值：120
	MaxVideoDuration int `json:"maxVideoDuration,omitempty"`

	// Onhold 是否在启动页面录制任务时暂停页面录制
	//
	// true：在启动页面录制任务时暂停页面录制。开启页面录制任务后立即暂停录制，录制服务会打开并渲染待录制页面，但不生成切片文件。
	//
	// false：启动页面录制任务并进行页面录制。
	//
	// 注意：建议你按照如下流程使用 onhold 字段：
	//
	// 调用 start 方法时将 onhold 设为 true，开启并暂停页面录制，自行判断页面录制开始的合适时机。
	//
	// 调用 update 并将 onhold 设为 false，继续进行页面录制。如果需要连续调用 update 方法暂停或继续页面录制，请在收到上一次 update 响应后再进行调用，否则可能导致请求结果与预期不一致。
	//
	// 默认值：false
	Onhold bool `json:"onhold"`

	// ReadyTimeout 设置页面加载超时时间，单位为秒
	//
	// 0 或不设置，表示不检测页面加载状态。
	//
	// [1,60] 之间的整数，表示页面加载超时时间。
	//
	// 0 <= ReadyTimeout <= 60
	//
	// 默认值：0
	ReadyTimeout int `json:"readyTimeout"`
}

type StartResp struct {
	Response
	SuccessResp StartSuccessResp
}

// StartSuccessResp 云端录制服务成功开始云端录制后返回的响应
type StartSuccessResp struct {
	// Cname 录制的频道名
	Cname string `json:"cname"`

	// UID 字符串内容为云端录制服务在 RTC 频道内使用的 UID，用于标识频道内的录制服务。
	UID string `json:"uid"`

	// ResourceId 云端录制资源 Resource ID。
	//
	// 使用这个 Resource ID 可以开始一段云端录制。这个 Resource ID 的有效期为 5 分钟，超时需要重新请求。
	ResourceId string `json:"resourceId"`

	// Sid 录制 ID。成功开始云端录制后，你会得到一个 Sid （录制 ID）。该 ID 是一次录制周期的唯一标识。
	Sid string `json:"sid"`
}

func (s *Start) WithForwardRegion(prefix core.ForwardedReginPrefix) *Start {
	s.forwardedRegionPrefix = prefix

	return s
}

func (s *Start) Do(ctx context.Context, resourceID string, mode string, payload *StartReqBody) (*StartResp, error) {
	path := s.buildPath(resourceID, mode)

	responseData, err := s.doRESTWithRetry(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp StartResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResp StartSuccessResp
		if err = responseData.UnmarshalToTarget(&successResp); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResp
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

const retryCount = 3

func (s *Start) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*core.BaseResponse, error) {
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
