package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Start struct {
	baseHandler
}

func NewStart(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Start {
	return &Start{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			retryCount: retryCount,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

const (
	IndividualMode = "individual"
	MixMode        = "mix"
	WebMode        = "web"
)

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/mode/{mode}/start
func (s *Start) buildPath(resourceID string, mode string) string {
	return s.prefixPath + "/resourceid/" + resourceID + "/mode/" + mode + "/start"
}

type StartReqBody struct {
	Cname         string              `json:"cname"`
	Uid           string              `json:"uid"`
	ClientRequest *StartClientRequest `json:"clientRequest"`
}

type StartClientRequest struct {
	Token                  string                  `json:"token,omitempty"`
	RecordingConfig        *RecordingConfig        `json:"recordingConfig,omitempty"`
	RecordingFileConfig    *RecordingFileConfig    `json:"recordingFileConfig,omitempty"`
	SnapshotConfig         *SnapshotConfig         `json:"snapshotConfig,omitempty"`
	StorageConfig          *StorageConfig          `json:"storageConfig,omitempty"`
	ExtensionServiceConfig *ExtensionServiceConfig `json:"extensionServiceConfig,omitempty"`
}

// @brief Configuration for recorded audio and video streams.
//
// @since v0.8.0
type RecordingConfig struct {
	// The channel type.(Required)
	//
	// The channel type can be set to:
	//
	//  - 0: The communication use-case (Default)
	//  - 1: Live streaming scene
	ChannelType int `json:"channelType"`

	// Subscribed media stream type.(Optional)
	//
	// The stream type can be set to:
	//  - 0: Subscribes to audio streams only. Suitable for smart voice review use-cases.
	//  - 1: Subscribes to video streams only.
	//  - 2: Subscribes to both audio and video streams.(Default)
	StreamTypes int `json:"streamTypes"`

	// Output mode of media stream.(Optional)
	//
	// The stream mode can be set to:
	//
	//  - "default": Default mode.
	//    Recording with audio transcoding will separately generate an M3U8 audio index file and a video index file.
	//  - "standard": Standard mode. Agora recommends using this mode.
	//    Recording with audio transcoding will separately generate an M3U8 audio index file, a video index file,
	//    and a merged audio and video index file. If VP8 encoding is used on the Web client, a merged MPD audio-video index file will be generated.
	//  - "original": Original encoding mode. It is applicable to individual non-transcoding audio recording.
	//    This field only takes effect when subscribing to audio only (streamTypes is 0).
	//    During the recording process, the audio is not transcoded, and an M3U8 audio index file is generated.
	StreamMode string `json:"streamMode,omitempty"`
	// The decryption mode.(Optional)
	//
	// If you have set channel encryption in the SDK client,
	// you need to set the same decryption mode for the cloud recording service.
	//
	// The decryption mode can be set to:
	//
	//  - 0: Not encrypted.(Default)
	//  - 1: AES_128_XTS encryption mode. 128-bit AES encryption, XTS mode.
	//  - 2: AES_128_ECB encryption mode. 128-bit AES encryption, ECB mode.
	//  - 3: AES_256_XTS encryption mode. 256-bit AES encryption, XTS mode.
	//  - 4: SM4_128_ECB encryption mode. 128-bit SM4 encryption, ECB mode.
	//  - 5: AES_128_GCM encryption mode. 128-bit AES encryption, GCM mode.
	//  - 6: AES_256_GCM encryption mode. 256-bit AES encryption, GCM mode.
	//  - 7: AES_128_GCM2 encryption mode. 128-bit AES encryption, GCM mode.
	//       Compared to AES_128_GCM encryption mode, AES_128_GCM2 encryption mode has higher security and requires setting a key and salt.
	//  - 8: AES_256_GCM2 encryption mode. 256-bit AES encryption, GCM mode.
	//       Compared to the AES_256_GCM encryption mode, the AES_256_GCM2 encryption mode is more secure and requires setting a key and salt.
	DecryptionMode int `json:"decryptionMode,omitempty"`
	// Keys related to encryption and decryption.(Optional)
	//
	// Only needs to be set when decryptionMode is not 0.
	Secret string `json:"secret,omitempty"`
	// Salt related to encryption and decryption.(Optional)
	//
	// Base64 encoding, 32-bit bytes.
	//
	// Only need to set when decryptionMode is 7 or 8.
	Salt string `json:"salt,omitempty"`
	// Set the sampling rate, bitrate, encoding mode, and number of channels for the output audio.(Optional)
	//
	// The audio profile can be set to:
	//
	//  - 0: 48 kHz sampling rate, music encoding, mono audio channel, and the encoding bitrate is about 48 Kbps.（Default）
	//  - 1: 48 kHz sampling rate, music encoding, mono audio channel, and the encoding bitrate is approximately 128 Kbps.
	//  - 2: 48 kHz sampling rate, music encoding, stereo audio channel, and the encoding bitrate is approximately 192 Kbps.
	AudioProfile int `json:"audioProfile,omitempty"`
	// Sets the stream type of the remote video.(Optional)
	//
	// If you enable dual-stream mode in the SDK client,
	// you can choose to subscribe to either the high-quality video stream or the low-quality video stream.
	//
	// The video stream type can be set to:
	//
	//  - 0: High-quality video stream refers to high-resolution and high-bitrate video stream.(Default)
	//  - 1: Low-quality video stream refers to low-resolution and low-bitrate video stream.
	VideoStreamType int `json:"videoStreamType,omitempty"`
	// Maximum channel idle time.(Optional)
	//
	// The unit is seconds.
	//
	// The value range is [5,259200].
	//
	// The default value is 30.
	MaxIdleTime int `json:"maxIdleTime,omitempty"`
	// Configurations for transcoded video output.(Optional)
	TranscodingConfig *TranscodingConfig `json:"transcodingConfig,omitempty"`
	// Specify which UIDs' audio streams to subscribe to.(Optional)
	//
	// If you want to subscribe to the audio stream of all UIDs, no need to set this field.
	SubscribeAudioUIDs []string `json:"subscribeAudioUids,omitempty"`
	// Specify which UIDs' audio streams not to subscribe to.(Optional)
	//
	// The cloud recording service will subscribe to the audio streams of all other UIDs except the specified ones.
	UnsubscribeAudioUIDs []string `json:"unSubscribeAudioUids,omitempty"`
	// Specify which UIDs' video streams to subscribe to.(Optional)
	//
	// If you want to subscribe to the video streams of all UIDs, no need to set this field.
	SubscribeVideoUIDs []string `json:"subscribeVideoUids,omitempty"`
	// Specify which UIDs' video streams not to subscribe to.(Optional)
	//
	// The cloud recording service will subscribe to the video streams of all UIDs except the specified ones.
	UnsubscribeVideoUIDs []string `json:"unSubscribeVideoUids,omitempty"`
	// Estimated peak number of subscribers.(Optional)
	//
	// The subscription group can be set to:
	//  - 0: 1 to 2 UIDs.
	//  - 1: 3 to 7 UIDs.
	//  - 2: 8 to 12 UIDs
	//  - 3: 13 to 17 UIDs
	//  - 4: 18 to 32 UIDs.
	//  - 5: 33 to 49 UIDs.
	SubscribeUidGroup int `json:"subscribeUidGroup,omitempty"`
}

// @brief Configurations for transcoded video output.
//
// @since v0.8.0
type TranscodingConfig struct {
	// The width of the video (pixels).(Optional)
	//
	// Width × Height cannot exceed 1920 × 1080.
	//
	// The default value is 360.
	Width int `json:"width,omitempty"`
	// The height of the video (pixels).(Optional)
	//
	// width × height cannot exceed 1920 × 1080.
	//
	// The default value is 640.
	Height int `json:"height,omitempty"`
	// The frame rate of the video (fps).(Optional)
	//
	// The default value is 15.
	FPS int `json:"fps,omitempty"`
	// The bitrate of the video (Kbps).(Optional)
	//
	// The default value is 1500.
	BitRate int `json:"bitrate,omitempty"`
	// Only need to set it in vertical layout.(Optional)
	//
	// Specify the user ID of the large video window.
	MaxResolutionUid string `json:"maxResolutionUid,omitempty"`
	// Composite video layout.(Optional)
	//
	// The video layout can be set to:
	//
	// - 0: Floating layout(Default).
	//   The first user to join the channel will be displayed as a large window, filling the entire canvas.
	//   The video windows of other users will be displayed as small windows, arranged horizontally from bottom to top,
	//   up to 4 rows, each with 4 windows. It supports up to a total of 17 windows of different users' videos.
	// - 1: Adaptive layout.
	//   Automatically adjust the size of each user's video window according to the number of users,
	//   each user's video window size is consistent, and supports up to 17 windows.
	// - 2: Vertical layout.
	//   The maxResolutionUid is specified to display the large video window on the left side of the screen,
	//   and the small video windows of other users are vertically arranged on the right side,
	//   with a maximum of two columns, 8 windows per column, supporting up to 17 windows.
	// - 3: Customized layout.
	//   Set the layoutConfig field to customize the mixed layout.
	MixedVideoLayout int `json:"mixedVideoLayout,omitempty"`
	// The background color of the video canvas.(Optional)
	//
	// The RGB color table is supported, with strings formatted as a # sign and 6 hexadecimal digits.
	//
	// The default value is "#000000", representing the black color.
	BackgroundColor string `json:"backgroundColor,omitempty"`
	// The URL of the background image of the video canvas.(Optional)
	//
	// The display mode of the background image is set to cropped mode.
	//
	// Cropped mode: Will prioritize to ensure that the screen is filled.
	// The background image size is scaled in equal proportion until the entire screen is filled with the background image.
	// If the length and width of the background image differ from the video window,
	// the background image will be peripherally cropped to fill the window.
	BackgroundImage string `json:"backgroundImage,omitempty"`
	// The URL of the default user screen background image.(Optional)
	DefaultUserBackgroundImage string `json:"defaultUserBackgroundImage,omitempty"`
	// Configurations of user's layout.(Optional)
	LayoutConfig []LayoutConfig `json:"layoutConfig,omitempty"`
	// Configurations of user's background image.(Optional)
	BackgroundConfig []BackgroundConfig `json:"backgroundConfig,omitempty"`
}

// @brief Configurations of user's layout.
//
// @since v0.8.0
type LayoutConfig struct {
	// The content of the string is the UID of the user to be displayed in the area, 32-bit unsigned integer.
	UID string `json:"uid"`
	// The relative value of the horizontal coordinate of the upper-left corner of the screen,
	// accurate to six decimal places.
	//
	// Layout from left to right, with 0.0 at the far left and 1.0 at the far right.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	XAxis float32 `json:"x_axis"`
	// The relative value of the vertical coordinate of the upper-left corner of this screen in the screen,
	// accurate to six decimal places.
	//
	// Layout from top to bottom, with 0.0 at the top and 1.0 at the bottom.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	YAxis float32 `json:"y_axis"`
	// The relative value of the width of this screen, accurate to six decimal places.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	Width float32 `json:"width"`
	// The relative value of the height of this screen, accurate to six decimal places.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	Height float32 `json:"height"`
	// The transparency of the user's video window. Accurate to six decimal places.
	//
	// 0.0 means the user's video window is transparent, and 1.0 indicates that it is completely opaque.
	//
	// The value range is [0,1].
	//
	// The default value is 1.
	Alpha float32 `json:"alpha,omitempty"`
	// The display mode of users' video windows.
	//
	// The rendering mode can be set to:
	//
	//  - 0: Cropped mode.(Default)
	//       Prioritize to ensure the screen is filled.
	//       The video window size is proportionally scaled until it fills the screen.
	//       If the video's length and width differ from the video window,
	//		 the video stream will be cropped from its edges to fit the window,
	//       under the aspect ratio set for the video window.
	//  - 1: Fit mode.
	//       Prioritize to ensure that all video content is displayed.
	//       The video size is scaled proportionally until one side of the video window is aligned with the screen border.
	//       If the video scale does not comply with the window size,
	//       the video will be scaled to fill the screen while maintaining its aspect ratio.
	//       This scaling may result in a black border around the edges of the video.
	RenderMode int `json:"render_mode"`
}

// @brief Configuration for the recorded files.
//
// @since v0.8.0
type RecordingFileConfig struct {
	// Type of video files generated by recording.(Optional)
	//
	// The file type can be set to:
	//
	//  - "hls": default value. M3U8 and TS files.
	//  - "mp4": MP4 files.
	AvFileType []string `json:"avFileType"`
}

// @brief Configuration for screenshot capture.
//
// @since v0.8.0
type SnapshotConfig struct {
	// The cycle for regular screenshots in the cloud recording.(Optional)
	//
	// The unit is seconds.
	//
	// The value range is [5,3600].
	//
	// The default value is 10.
	CaptureInterval int `json:"captureInterval,omitempty"`
	// The file format of screenshots.
	//
	// Currently only ["jpg"] is supported, which generates screenshot files in JPG format.
	FileType []string `json:"fileType"`
}

// @brief Configuration for third-party cloud storage.
//
// @since v0.8.0
type StorageConfig struct {
	// Third-party cloud storage platforms.(Required)
	//
	// The vendor can be set to:
	//
	//  - 1: Amazon S3
	//  - 2: Alibaba Cloud
	//  - 3: Tencent Cloud
	//  - 5: Microsoft Azure
	//  - 6: Google Cloud
	//  - 7: Huawei Cloud
	//  - 8: Baidu IntelligentCloud
	//  - 11: Self-built cloud storage
	Vendor int `json:"vendor"`

	// The region information specified for the third-party cloud storage.(Required)
	Region int `json:"region"`

	// Third-party cloud storage bucket.(Required)
	Bucket string `json:"bucket"`

	// The access key of third-party cloud storage.(Required)
	AccessKey string `json:"accessKey"`

	// A temporary security token for third-party cloud storage.
	// This token is issued by the cloud service provider's Security Token Service (STS) and used to grant limited access rights to third-party cloud storage resources.
	//
	// Currently supported cloud service providers include only the following:
	//
	//  - 1: Amazon S3
	//  - 2: Alibaba Cloud
	//  - 3: Tencent Cloud.
	StsToken string `json:"stsToken,omitempty"`

	// The stsToken expiration timestamp used to mark UNIX time, in seconds.(Optional)
	StsExpiration int `json:"stsExpiration,omitempty"`

	// The secret key of third-party cloud storage.(Required)
	SecretKey string `json:"secretKey"`

	// The storage location of the recorded files in the third-party cloud is related to the prefix of the file name.(Optional)
	FileNamePrefix []string `json:"fileNamePrefix,omitempty"`

	// Third-party cloud storage services will encrypt and tag the uploaded recording files according to this field.(Optional)
	ExtensionParams *ExtensionParams `json:"extensionParams,omitempty"`
}

// @brief Third-party cloud storage services will encrypt and tag the uploaded recording files according to this field.
//
// @since v0.8.0
type ExtensionParams struct {
	// The encryption mode.(Required)
	//
	// This field is only applicable to Amazon S3,
	// and the value can be set to:
	//
	//  - "kms": KMS encryption.
	//  - "aes256": AES256 encryption.
	SSE string `json:"sse"`
	// Tag content.(Required)
	//
	// After setting this field, the third-party cloud storage service
	// will tag the uploaded recording files according to the content of this tag.
	// This field is only applicable to Alibaba Cloud and Amazon S3.
	// For other third-party cloud storage services, this field is not required.
	Tag string `json:"tag"`
	// Domain name of self-built cloud storage.(Optional)
	//
	// This field is required when vendor is set to 11.
	Endpoint string `json:"endpoint,omitempty"`
}

// @brief Configurations for extended services.
//
// @since v0.8.0
type ExtensionServiceConfig struct {
	// Error handling policy.(Optional)
	//
	// You can only set it to the default value, "error_abort",
	// which means that once an error occurs to an extension service,
	// all other non-extension services, such as stream subscription, also stop.
	ErrorHandlePolicy string `json:"errorHandlePolicy,omitempty"`
	// Extended services.(Required)
	ExtensionServices []ExtensionService `json:"extensionServices"`
}

type ServiceParamInterface interface {
	ServiceParam()
}

// @brief Configuration for extended services.
//
// @since v0.8.0
type ExtensionService struct {
	// Name of the extended service.(Required)
	//
	// The service name can be set to:
	//
	//  - "web_recorder_service": Represents the extended service is web page recording.
	//  - "rtmp_publish_service": Represents the extended service is to push web page recording to the CDN.
	ServiceName string `json:"serviceName"`
	// Error handling strategy within the extension service.(Optional)
	//
	// The error handling strategy can be set to:
	//
	//  - "error_abort": the default and only value during web page recording.
	//    Stop other extension services when the current extension service encounters an error.
	//  - "error_ignore": The only default value when you push the web page recording to the CDN.
	//    Other extension services are not affected when the current extension service encounters an error.
	ErrorHandlePolicy string `json:"errorHandlePolicy"`
	// Specific configurations for extension services.(Required)
	//
	// - "WebRecordingServiceParam" for web page recording. See WebRecordingServiceParam for details.
	// - "RtmpPublishServiceParam" for pushing web page recording to the CDN. See RtmpPublishServiceParam for details.
	ServiceParam ServiceParamInterface `json:"serviceParam"`
}

// @brief The CDN address to which you push the stream.
//
// @since v0.8.0
type Outputs struct {
	// The CDN address to which you push the stream.(Required)
	RtmpURL string `json:"rtmpUrl"`
}

// @brief Service parameter configuration for web page recording.
//
// @since v0.8.0
type WebRecordingServiceParam struct {
	// The address of the page to be recorded.(Required)
	URL string `json:"url"`
	// The bitrate of the output video (Kbps).(Optional)
	//
	// For different output video resolutions, the default value of videoBitrate is different:
	//
	//  - Output video resolution is greater than or equal to 1280 × 720, and the default value is 2000.
	//  - Output video resolution is less than 1280 × 720, and the default value is 1500.
	VideoBitRate int `json:"VideoBitrate,omitempty"`
	// The frame rate of the output video (fps).(Optional)
	//
	// The value range is [5,60].
	//
	// The default value is 15.
	VideoFPS int `json:"videoFps,omitempty"`
	// Sampling rate, bitrate, encoding mode, and number of channels for the audio output.(Required)
	//
	// The audio profile can be set to:
	//
	//  - 0: 48 kHz sampling rate, music encoding, mono audio channel, and the encoding bitrate is approximately 48 Kbps.
	//  - 1: 48 kHz sampling rate, music encoding, mono audio channel, and the encoding bitrate is approximately 128 Kbps.
	//  - 2: 48 kHz sampling rate, music encoding, stereo audio channel, and the encoding bitrate is approximately 192 Kbps.
	AudioProfile int `json:"audioProfile"`
	// Whether to enable the mobile web mode.(Optional)
	//
	//  - true: Enables the mode. After enabling,
	// 	  the recording service uses the mobile web rendering mode to record the current page.
	//  - false: Disables the mode.(Default)
	Mobile bool `json:"mobile,omitempty"`
	// The output video width (pixel).(Required)
	//
	// The product of videoWidth and videoHeight should be less than or equal to 1920 × 1080.
	VideoWidth int `json:"videoWidth"`
	// The output video height (pixel).(Required)
	//
	// The product of videoWidth and videoHeight should be less than or equal to 1920 × 1080.
	VideoHeight int `json:"videoHeight"`
	// The maximum duration of web page recording (hours). (Required)
	//
	// The web page recording will automatically stop after exceeding this value.
	//
	// The value range is [1,720].
	MaxRecordingHour int `json:"maxRecordingHour"`
	// Maximum length of MP4 slice file generated by web page recording, in minutes.(Optional)
	//
	// During the web page recording process,
	// the recording service will create a new MP4 slice file when the current MP4 file
	// duration exceeds the maxVideoDuration approximately.
	//
	// The value range is [30,240].
	//
	// The default value is 120.
	MaxVideoDuration int `json:"maxVideoDuration,omitempty"`
	// Whether to pause page recording when starting a web page recording task.(Optional)
	//
	// - true: Pauses the web page recording that has been started.
	//   Immediately pause the recording after starting the web page recording task.
	//   The recording service will open and render the page to be recorded, but will not generate slice files.
	// - false: Starts a web page recording task and performs web page recording.(Default)
	Onhold bool `json:"onhold"`
	// Set the page load timeout in seconds.(Optional)
	//
	// The value range is [0,60].
	//
	// The default value is 0.
	ReadyTimeout int `json:"readyTimeout"`
}

func (w *WebRecordingServiceParam) ServiceParam() {
}

// @brief Service parameter configuration for pushing web page recording to the CDN.
//
// @since v0.8.0
type RtmpPublishServiceParam struct {
	// The array of CDN addresses to which you push the stream.(Required)
	Outputs []Outputs `json:"outputs,omitempty"`
}

func (r *RtmpPublishServiceParam) ServiceParam() {
}

// @brief StartResp returned by the various of cloud recording scenarios Start API.
//
// @since v0.8.0
type StartResp struct {
	// Response returned by the cloud recording API, see Response for details
	Response
	// Successful response, see StartSuccessResp for details
	SuccessResponse StartSuccessResp
}

// @brief Successful response returned by the various of cloud recording scenarios Start API.
//
// @since v0.8.0
type StartSuccessResp struct {
	// Channel name
	Cname string `json:"cname"`
	// User ID
	UID string `json:"uid"`
	// Unique identifier of the resource
	ResourceId string `json:"resourceId"`
	// Unique identifier of the recording session
	Sid string `json:"sid"`
}

func (s *Start) Do(ctx context.Context, resourceID string, mode string, payload *StartReqBody) (*StartResp, error) {
	path := s.buildPath(resourceID, mode)

	responseData, err := doRESTWithRetry(ctx, s.module, s.logger, s.retryCount, s.client, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
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
		resp.SuccessResponse = successResp
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
