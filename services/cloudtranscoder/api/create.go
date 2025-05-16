package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Create struct {
	baseHandler
}

func NewCreate(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Create {
	return &Create{
		baseHandler: baseHandler{
			module: module, logger: logger, retryCount: retryCount, client: client, prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks?builderToken={tokenName}
func (c *Create) buildPath(tokenName string) string {
	return c.prefixPath + "/tasks?builderToken=" + tokenName
}

type CreateReqBody struct {
	Services *CreateReqServices `json:"services,omitempty"`
}

type CreateReqServices struct {
	CloudTranscoder *CloudTranscoderPayload `json:"cloudTranscoder,omitempty"`
}

type CloudTranscoderPayload struct {
	// Service type, here it is "cloudTranscoderV2"
	ServiceType string                 `json:"serviceType,omitempty"`
	Config      *CloudTranscoderConfig `json:"config,omitempty"`
}

type CloudTranscoderConfig struct {
	Transcoder *CloudTranscoderConfigPayload `json:"transcoder,omitempty"`
}

type CloudTranscoderConfigPayload struct {
	// Maximum idle time (in seconds) for cloud transcoder. Idle means all broadcasters corresponding to the audio/video streams
	// processed by the cloud transcoder have left the channel.
	// After being idle for the set idleTimeOut, the cloud transcoder will be automatically destroyed.
	//
	// Range: [1,86400]
	//
	// Default: 300
	IdleTimeout uint                        `json:"idleTimeout"`
	AudioInputs []CloudTranscoderAudioInput `json:"audioInputs,omitempty"`
	VideoInputs []CloudTranscoderVideoInput `json:"videoInputs,omitempty"`
	Canvas      *CloudTranscoderCanvas      `json:"canvas,omitempty"`
	WaterMarks  []CloudTranscoderWaterMark  `json:"waterMarks,omitempty"`
	Outputs     []CloudTranscoderOutput     `json:"outputs,omitempty"`
}

type CloudTranscoderAudioInput struct {
	Rtc *CloudTranscoderRtc `json:"rtc,omitempty"`
}

type CloudTranscoderRtc struct {
	// RTC channel name for the audio/video input source (or output)
	//
	// Currently only supports subscribing to audio/video sources from a single channel.
	// Audio and video sources must belong to the same channel.
	RtcChannel string `json:"rtcChannel,omitempty"`

	// UID corresponding to the audio/video input source (or output)
	//
	// Duplicate UIDs are not allowed in an RTC channel, so ensure this value differs from other users' UIDs in the channel.
	//
	// Note: The default value of rtcUid is 0, indicating that AudioInputs will use full-channel audio mixing
	RtcUID int `json:"rtcUid"`

	// Token required for the cloud transcoder to enter the RTC channel of the video source to be transcoded (or the transcoded output audio/video stream).
	//
	// This value can be used to ensure channel security and prevent unauthorized users from disrupting other users in the channel.
	//
	// Note:
	//   - When configuring input streams, the UID of the cloud transcoder in the RTC channel is randomly assigned by Agora.
	//     Therefore, when generating the Token, you must use uid=0.
	//   - When configuring output streams, the UID of the cloud transcoder in the RTC channel is the outputs.rtc.rtcUid you specified.
	//     Therefore, when generating the Token, you must use the same uid as outputs.rtc.rtcUid
	RtcToken string `json:"rtcToken,omitempty"`
}

type CloudTranscoderVideoInput struct {
	Rtc *CloudTranscoderRtc `json:"rtc,omitempty"`

	// URL of the placeholder image when the user is offline.
	//
	// Must be a valid URL and include a jpg or png suffix.
	PlaceholderImageURL string                 `json:"placeholderImageUrl,omitempty"`
	Region              *CloudTranscoderRegion `json:"region,omitempty"`
}

type CloudTranscoderRegion struct {
	// X coordinate (px) of the image on the canvas.
	//
	// With the top-left corner of the canvas as the origin, the x-coordinate is the horizontal displacement of the top-left corner of the image relative to the origin.
	//
	// Range: [0,120]
	X uint `json:"x"`

	// Y coordinate (px) of the image on the canvas.
	//
	// With the top-left corner of the canvas as the origin, the y-coordinate is the vertical displacement of the top-left corner of the image relative to the origin.
	//
	// Range: [0,3840]
	Y uint `json:"y"`

	// Width of the image (px).
	//
	// Range: [120,3840]
	Width uint `json:"width,omitempty"`

	// Height of the image (px).
	//
	// Range: [120,3840]
	Height uint `json:"height,omitempty"`

	// Layer number of the image.
	//  - 0 represents the bottom layer.
	//  - 100 represents the top layer.
	//
	// Range: [0,100]
	ZOrder uint `json:"zOrder"`
}

type CloudTranscoderCanvas struct {
	// Width of the canvas (px).
	//
	// Range: [120,3840]
	Width uint `json:"width,omitempty"`

	// Height of the canvas (px).
	//
	// Range: [120,3840]
	Height uint `json:"height,omitempty"`

	// Background color of the canvas.
	//
	// RGB color value, expressed as a decimal number.
	//
	// For example, 0 represents black, 255 represents blue.
	//
	// Range: [0,16777215]
	Color uint `json:"color"`

	// Background image of the canvas.
	//
	// Must be a valid URL and include a jpg or png suffix.
	//
	// Note: If no value is provided, there will be no canvas background image.
	BackgroundImage string `json:"backgroundImage,omitempty"`

	// Fill mode for the canvas background image:
	//
	//  - "FILL": Scale the image while maintaining the aspect ratio, and crop it centered.
	//  - "FIT": Scale the image while maintaining the aspect ratio, ensuring it is fully displayed.
	//
	// Default: "FILL"
	FillMode string `json:"fillMode,omitempty"`
}

type CloudTranscoderWaterMark struct {
	// URL of the watermark image.
	//
	// Must be a valid URL and include a jpg or png suffix.
	ImageURL string                 `json:"imageUrl,omitempty"`
	Region   *CloudTranscoderRegion `json:"region,omitempty"`

	// Fill mode for the canvas background image:
	//
	//  - "FILL": Scale the image while maintaining the aspect ratio, and crop it centered.
	//  - "FIT": Scale the image while maintaining the aspect ratio, ensuring it is fully displayed.
	//
	// Default: "FILL"
	FillMode string `json:"fillMode,omitempty"`
}

type CloudTranscoderOutputAudioOption struct {
	// Audio properties for transcoding output:
	//   - "AUDIO_PROFILE_DEFAULT": 48 kHz sampling rate, music encoding, mono, maximum encoding bitrate of 64 Kbps.
	//   - "AUDIO_PROFILE_SPEECH_STANDARD": 32 kHz sampling rate, speech encoding, mono, maximum encoding bitrate of 18 Kbps.
	//   - "AUDIO_PROFILE_MUSIC_STANDARD": 48 KHz sampling rate, music encoding, mono, maximum encoding bitrate of 64 Kbps.
	//   - "AUDIO_PROFILE_MUSIC_STANDARD_STEREO": 48 KHz sampling rate, music encoding, stereo, maximum encoding bitrate of 80 Kbps.
	//   - "AUDIO_PROFILE_MUSIC_HIGH_QUALITY": 48 KHz sampling rate, music encoding, mono, maximum encoding bitrate of 96 Kbps.
	// 	 - "AUDIO_PROFILE_MUSIC_HIGH_QUALITY_STEREO": 48 KHz sampling rate, music encoding, stereo, maximum encoding bitrate of 128 Kbps.
	//
	// Default: "AUDIO_PROFILE_DEFAULT"
	ProfileType string `json:"profileType,omitempty"`
}

type CloudTranscoderOutputVideoOption struct {
	// Frame rate (fps) of the transcoded output video.
	//
	// Range: [1,30]
	//
	// Default: 15
	FPS uint `json:"fps,omitempty"`

	// Codec for the transcoded output video. Values include:
	//  - "H264": Standard H.264 encoding.
	//  - "VP8": Standard VP8 encoding.
	Codec string `json:"codec,omitempty"`

	// Bitrate of the transcoded output video.
	//
	// Range: [1,10000]
	//
	// Note: If you do not provide a value, Agora will automatically set the video bitrate based on network conditions and other video properties.
	Bitrate uint `json:"bitrate,omitempty"`

	// Width of the image (px).
	//
	// Range: [120,3840]
	Width uint `json:"width,omitempty"`

	// Height of the image (px).
	//
	// Range: [120,3840]
	Height uint `json:"height,omitempty"`
}

type CloudTranscoderOutput struct {
	Rtc         *CloudTranscoderRtc               `json:"rtc,omitempty"`
	AudioOption *CloudTranscoderOutputAudioOption `json:"audioOption,omitempty"`
	VideoOption *CloudTranscoderOutputVideoOption `json:"videoOption,omitempty"`
}

type CreateSuccessResp struct {
	// Transcoding task ID, a UUID used to identify the cloud transcoder for this request operation
	TaskID string `json:"taskId"`

	// Unix timestamp (seconds) when the transcoding task was created
	CreateTs int64 `json:"createTs"`

	// Running status of the transcoding task:
	//  - "IDLE": Task has not started.
	//
	//  - "PREPARED": Task has received a start request.
	//
	//  - "STARTING": Task is starting.
	//
	//  - "CREATED": Task initialization complete.
	//
	//  - "STARTED": Task has started.
	//
	//  - "IN_PROGRESS": Task is in progress.
	//
	//  - "STOPPING": Task is stopping.
	//
	//  - "STOPPED": Task has stopped.
	//
	//  - "EXIT": Task exited normally.
	//
	//  - "FAILURE_STOP": Task exited abnormally.
	//
	// Note: You can use this field to monitor the task status.
	Status string `json:"status"`
}

type CreateResp struct {
	Response
	SuccessResp CreateSuccessResp
}

func (c *Create) Do(ctx context.Context, tokenName string, payload *CreateReqBody) (*CreateResp, error) {
	path := c.buildPath(tokenName)
	responseData, err := doRESTWithRetry(ctx, c.module, c.logger, c.retryCount, c.client, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
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
		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		resp.ErrResponse = errResponse
	}

	resp.BaseResponse = responseData

	return &resp, nil
}
