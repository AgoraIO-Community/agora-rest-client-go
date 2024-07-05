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
	Cloudtranscoder *CloudtranscoderPayload `json:"cloudTranscoder"`
}

type CloudtranscoderPayload struct {
	ServiceType string                 `json:"serviceType"`
	Config      *CloudTranscoderConfig `json:"config"`
}

type CloudTranscoderConfig struct {
	Transcoder *CloudTranscoderConfigPayload `json:"transcoder"`
}

type CloudTranscoderConfigPayload struct {
	IdleTimeout uint                        `json:"idleTimeout"`
	AuidoInputs []CloudTranscoderAudioInput `json:"audioInputs"`
	VideoInputs []CloudTranscoderVideoInput `json:"videoInputs"`
	Cavas       *CloudTranscoderCanvas      `json:"canvas"`
	WaterMarks  []CloudTranscoderWaterMark  `json:"waterMarks"`
	Outputs     []CloudTranscoderOutput     `json:"outputs"`
}

type CloudTranscoderAudioInput struct {
	Rtc *CloudTranscoderRtcInput `json:"rtc"`
}

type CloudTranscoderRtcInput struct {
	RtcChannel string `json:"rtcChannel"`
	RtcUID     int    `json:"rtcUid"`
	RtcToken   string `json:"rtcToken"`
}

type CloudTranscoderVideoInput struct {
	Rtc                 *CloudTranscoderRtcInput `json:"rtc"`
	PlaceholderImageURL string                   `json:"placeholderImageUrl"`
	Region              *CloudTranscoderRegion   `json:"region"`
}

type CloudTranscoderRegion struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	ZOrder int     `json:"zOrder"`
}

type CloudTranscoderCanvas struct {
	Width           uint   `json:"width"`
	Height          uint   `json:"height"`
	Color           uint   `json:"color"`
	BackgroundImage string `json:"backgroundImage"`
	FillMode        string `json:"fillMode"`
}

type CloudTranscoderWaterMark struct {
	ImageURL string                 `json:"imageUrl"`
	Region   *CloudTranscoderRegion `json:"region"`
	FillMode string                 `json:"fillMode"`
}

type CloudTranscoderOutputAduioOption struct {
	ProfileType string `json:"profileType"`
}

type CloudTranscoderOutputVideoOption struct {
	FPS                   uint   `json:"fps"`
	Codec                 string `json:"codec"`
	Bitrate               uint   `json:"bitrate"`
	Width                 uint   `json:"width"`
	Height                uint   `json:"height"`
	LowBitrateHighQuality bool   `json:"lowBitrateHighQuality"`
}

type CloudTranscoderOutput struct {
	Rtc         *CloudTranscoderRtcInput          `json:"rtc"`
	AudioOption *CloudTranscoderOutputAduioOption `json:"audioOption"`
	VideoOption *CloudTranscoderOutputVideoOption `json:"videoOption"`
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
