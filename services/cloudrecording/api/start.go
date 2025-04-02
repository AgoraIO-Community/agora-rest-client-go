package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
)

type Start struct {
	module     string
	logger     log.Logger
	client     client.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

func NewStart(module string, logger log.Logger, client client.Client, prefixPath string) *Start {
	return &Start{module: module, logger: logger, client: client, prefixPath: prefixPath}
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

type AppsCollection struct {
	CombinationPolicy string `json:"combinationPolicy,omitempty"`
}

type RecordingConfig struct {
	ChannelType          int                `json:"channelType"`
	StreamTypes          int                `json:"streamTypes"`
	StreamMode           string             `json:"streamMode,omitempty"`
	DecryptionMode       int                `json:"decryptionMode,omitempty"`
	Secret               string             `json:"secret,omitempty"`
	Salt                 string             `json:"salt,omitempty"`
	AudioProfile         int                `json:"audioProfile,omitempty"`
	VideoStreamType      int                `json:"videoStreamType,omitempty"`
	MaxIdleTime          int                `json:"maxIdleTime,omitempty"`
	TranscodingConfig    *TranscodingConfig `json:"transcodingConfig,omitempty"`
	SubscribeAudioUIDs   []string           `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUIDs []string           `json:"unSubscribeAudioUids,omitempty"`
	SubscribeVideoUIDs   []string           `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUIDs []string           `json:"unSubscribeVideoUids,omitempty"`
	SubscribeUidGroup    int                `json:"subscribeUidGroup"`
}

type Container struct {
	Format string `json:"format"`
}

type TranscodeOptions struct {
	Container   *Container   `json:"container"`
	TransConfig *TransConfig `json:"transConfig"`
	Audio       *Audio       `json:"audio,omitempty"`
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
	Width                      int                `json:"width,omitempty"`
	Height                     int                `json:"height,omitempty"`
	FPS                        int                `json:"fps,omitempty"`
	BitRate                    int                `json:"bitrate,omitempty"`
	MaxResolutionUid           string             `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout           int                `json:"mixedVideoLayout,omitempty"`
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
	Alpha      float32 `json:"alpha,omitempty"`
	RenderMode int     `json:"render_mode"`
}

type RecordingFileConfig struct {
	AvFileType []string `json:"avFileType"`
}

type SnapshotConfig struct {
	CaptureInterval int      `json:"captureInterval,omitempty"`
	FileType        []string `json:"fileType"`
}

type StorageConfig struct {
	Vendor          int              `json:"vendor"`
	Region          int              `json:"region"`
	Bucket          string           `json:"bucket"`
	AccessKey       string           `json:"accessKey"`
	SecretKey       string           `json:"secretKey"`
	FileNamePrefix  []string         `json:"fileNamePrefix,omitempty"`
	ExtensionParams *ExtensionParams `json:"extensionParams,omitempty"`
}

type ExtensionParams struct {
	SSE string `json:"sse"`
	Tag string `json:"tag"`
}

type ExtensionServiceConfig struct {
	ErrorHandlePolicy string             `json:"errorHandlePolicy,omitempty"`
	ExtensionServices []ExtensionService `json:"extensionServices"`
}

type ServiceParamInterface interface {
	ServiceParam()
}

type ExtensionService struct {
	ServiceName       string               `json:"serviceName"`
	ErrorHandlePolicy string               `json:"errorHandlePolicy"`
	ServiceParam      ServiceParamInterface `json:"serviceParam"`
}

type Outputs struct {
	RtmpURL string `json:"rtmpUrl"`
}

type WebRecordingServiceParam struct {
	URL              string `json:"url"`
	VideoBitRate     int    `json:"VideoBitrate,omitempty"`
	VideoFPS         int    `json:"videoFps,omitempty"`
	AudioProfile     int    `json:"audioProfile"`
	Mobile           bool   `json:"mobile,omitempty"`
	VideoWidth       int    `json:"videoWidth"`
	VideoHeight      int    `json:"videoHeight"`
	MaxRecordingHour int    `json:"maxRecordingHour"`
	MaxVideoDuration int    `json:"maxVideoDuration,omitempty"`
	Onhold           bool   `json:"onhold"`
	ReadyTimeout     int    `json:"readyTimeout"`
}

func (w *WebRecordingServiceParam) ServiceParam() {
}

type RtmpPublishServiceParam struct {
	Outputs []Outputs `json:"outputs,omitempty"`
}

func (r *RtmpPublishServiceParam) ServiceParam() {
}

type StartResp struct {
	Response
	SuccessResponse StartSuccessResp
}

type StartSuccessResp struct {
	Cname      string `json:"cname"`
	UID        string `json:"uid"`
	ResourceId string `json:"resourceId"`
	Sid        string `json:"sid"`
}

func (s *Start) Do(ctx context.Context, resourceID string, mode string, payload *StartReqBody) (*StartResp, error) {
	path := s.buildPath(resourceID, mode)

	responseData, err := s.doRESTWithRetry(ctx, path, http.MethodPost, payload)
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

const startRetryCount = 3

func (s *Start) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		resp       *agora.BaseResponse
		err        error
		retryCount int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		resp, doErr = s.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			s.logger.Debugf(ctx, s.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
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
		return retryCount >= startRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		s.logger.Debugf(ctx, s.module, "http request err:%s", err)
		retryCount++
	})

	return resp, err
}
