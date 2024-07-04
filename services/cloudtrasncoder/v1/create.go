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
	CloudTrasncoder *CloudTrasncoderPayload `json:"cloudTranscoder"`
}

type CloudTrasncoderPayload struct {
	ServiceType string                 `json:"serviceType"`
	Config      *CloudTrasncoderConfig `json:"config"`
}

type CloudTrasncoderConfig struct {
	Transcoder *CloudTrasncoderConfigPayload `json:"transcoder"`
}

type CloudTrasncoderConfigPayload struct {
	IdleTimeout uint                        `json:"idleTimeout"`
	AuidoInputs []CloudTrasncoderAudioInput `json:"audioInputs"`
	VideoInputs []CloudTrasncoderVideoInput `json:"videoInputs"`
	Cavas       *CloudTrasncoderCanvas      `json:"canvas"`
	WaterMarks  []CloudTranscoderWaterMark  `json:"waterMarks"`
	Outputs     []CloudTrasncoderOutput     `json:"outputs"`
}

type CloudTrasncoderAudioInput struct {
	Rtc *CloudTrasncoderRtcInput `json:"rtc"`
}

type CloudTrasncoderRtcInput struct {
	RtcChannel string `json:"rtcChannel"`
	RtcUID     int    `json:"rtcUid"`
	RtcToken   string `json:"rtcToken"`
}

type CloudTrasncoderVideoInput struct {
	Rtc                 *CloudTrasncoderRtcInput `json:"rtc"`
	PlaceholderImageURL string                   `json:"placeholderImageUrl"`
	Region              *CloudTrasncoderRegion   `json:"region"`
}

type CloudTrasncoderRegion struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	ZOrder int     `json:"zOrder"`
}

type CloudTrasncoderCanvas struct {
	Width           uint   `json:"width"`
	Height          uint   `json:"height"`
	Color           uint   `json:"color"`
	BackgroundImage string `json:"backgroundImage"`
	FillMode        string `json:"fillMode"`
}

type CloudTranscoderWaterMark struct {
	ImageURL string                 `json:"imageUrl"`
	Region   *CloudTrasncoderRegion `json:"region"`
	FillMode string                 `json:"fillMode"`
}

type CloudTrasncoderOutputAduioOption struct {
	ProfileType string `json:"profileType"`
}

type CloudTrasncoderOutputVideoOption struct {
	FPS                   uint   `json:"fps"`
	Codec                 string `json:"codec"`
	Bitrate               uint   `json:"bitrate"`
	Width                 uint   `json:"width"`
	Height                uint   `json:"height"`
	LowBitrateHighQuality bool   `json:"lowBitrateHighQuality"`
}

type CloudTrasncoderOutput struct {
	Rtc         *CloudTrasncoderRtcInput          `json:"rtc"`
	AudioOption *CloudTrasncoderOutputAduioOption `json:"audioOption"`
	VideoOption *CloudTrasncoderOutputVideoOption `json:"videoOption"`
}

type CreateSuccessResp struct {
	TaskID   string `json:"taskId"`
	CreateTs int64  `json:"createTs"`
	Status   string `json:"status"`
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
