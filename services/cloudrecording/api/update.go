package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Update struct {
	baseHandler
}

func NewUpdate(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Update {
	return &Update{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			retryCount: retryCount,
			client:     client,
			prefixPath: prefixPath,
		},
	}
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

// @brief Update subscription lists.
//
// @since v0.8.0
type UpdateStreamSubscribe struct {
	// The audio subscription list.(Optional)
	AudioUidList *UpdateAudioUIDList `json:"audioUidList,omitempty"`

	// The video subscription list.(Optional)
	VideoUidList *UpdateVideoUIDList `json:"videoUidList,omitempty"`
}

// @brief Update audio subscription list.
//
// @since v0.8.0
type UpdateAudioUIDList struct {
	// Specify which UIDs' audio streams to subscribe to.
	//
	// If you want to subscribe to the audio stream of all UIDs, no need to set this field.
	//
	// Set as ["#allstream#"] to subscribe to the audio streams of all UIDs in the channel.
	SubscribeAudioUIDs []string `json:"subscribeAudioUids,omitempty"`

	// Specify which UIDs' audio streams not to subscribe to.
	//
	// The cloud recording service will subscribe to the audio streams of all other UIDs except the specified ones.
	UnsubscribeAudioUIDs []string `json:"unSubscribeAudioUids,omitempty"`
}

// @brief Update video subscription list.
//
// @since v0.8.0
type UpdateVideoUIDList struct {
	// Specify which UIDs' video streams to subscribe to.
	//
	// If you want to subscribe to the video stream of all UIDs, no need to set this field.
	//
	// Set as ["#allstream#"] to subscribe to the video streams of all UIDs in the channel.
	SubscribeVideoUIDs []string `json:"subscribeVideoUids,omitempty"`

	// Specify which UIDs' video streams not to subscribe to.
	//
	// The cloud recording service will subscribe to the video streams of all other UIDs except the specified ones.
	UnsubscribeVideoUIDs []string `json:"unSubscribeVideoUids,omitempty"`
}

// @brief Used to update the web page recording configurations
//
// @since v0.8.0
type UpdateWebRecordingConfig struct {
	// Set whether to pause the web page recording.(Optional)
	//
	//  - true: Pauses web page recording and generating recording files.
	//  - false: Continues web page recording and generates recording files.(Default)
	Onhold bool `json:"onhold"`
}

// @brief Used to update the configurations for pushing web page recording to the CDN.
//
// @since v0.8.0
type UpdateRtmpPublishConfig struct {
	// The array of CDN URL where you push the stream to.(Optional)
	Outputs []UpdateOutput `json:"outputs"`
}

// @brief The CDN URL where you push the stream to.
//
// @since v0.8.0
type UpdateOutput struct {
	// The CDN URL where you push the stream to.(Optional)
	RtmpURL string `json:"rtmpUrl"`
}

// @brief UpdateResp returned by the various of cloud recording scenarios Update API.
type UpdateResp struct {
	// Response returned by the cloud recording API, see Response for details
	Response
	// Successful response, see UpdateSuccessResp for details
	SuccessResponse UpdateSuccessResp
}

// @brief Successful response returned by the various of cloud recording scenarios Update API.
type UpdateSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string `json:"resourceId"`
	// Unique identifier of the recording session
	Sid string `json:"sid"`
	// User ID used by the cloud recording service in the RTC channel to identify the recording service in the channel
	UID string `json:"uid"`
	// Name of the channel to be recorded
	Cname string `json:"cname"`
}

func (u *Update) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateReqBody) (*UpdateResp, error) {
	path := u.buildPath(resourceID, sid, mode)

	responseData, err := doRESTWithRetry(ctx, u.module, u.logger, u.retryCount, u.client, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
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
