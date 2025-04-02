package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type Update struct {
	client     client.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

func NewUpdate(client client.Client, prefixPath string) *Update {
	return &Update{client: client, prefixPath: prefixPath}
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

type UpdateStreamSubscribe struct {
	AudioUidList *UpdateAudioUIDList `json:"audioUidList,omitempty"`
	VideoUidList *UpdateVideoUIDList `json:"videoUidList,omitempty"`
}

type UpdateAudioUIDList struct {
	SubscribeAudioUIDs    []string `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUIDs []string `json:"unSubscribeAudioUids,omitempty"`
}

type UpdateVideoUIDList struct {
	SubscribeVideoUIDs    []string `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUIDs []string `json:"unSubscribeVideoUids,omitempty"`
}

type UpdateWebRecordingConfig struct {
	Onhold bool `json:"onhold"`
}

type UpdateRtmpPublishConfig struct {
	Outputs []UpdateOutput `json:"outputs"`
}

type UpdateOutput struct {
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
