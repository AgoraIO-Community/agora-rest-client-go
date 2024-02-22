package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Update struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/update
func (s *Update) buildPath(resourceID string, sid string, mode string) string {
	return s.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/update"
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

type UpdateWebRecordingClientRequest struct {
	WebRecordingConfig *UpdateWebRecordingConfig `json:"webRecordingConfig,omitempty"`
	RtmpPublishConfig  *UpdateRtmpPublishConfig  `json:"rtmpPublishConfig,omitempty"`
}

type UpdateStreamSubscribe struct {
	AudioUidList *UpdateAudioUIDList `json:"audioUidList,omitempty"`
	VideoUidList *UpdateVideoUIDList `json:"videoUidList,omitempty"`
}

type UpdateAudioUIDList struct {
	SubscribeAudioUIDs   []string `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUIDs []string `json:"unsubscribeAudioUids,omitempty"`
}

type UpdateVideoUIDList struct {
	SubscribeVideoUIDs   []string `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUIDs []string `json:"unsubscribeVideoUids,omitempty"`
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
	SuccessResp UpdateSuccessResp
}

type UpdateSuccessResp struct {
	ResourceId string `json:"resourceId"`
	SID        string `json:"sid"`
	UID        string `json:"uid"`
	Cname      string `json:"cname"`
}

type UpdateServerResponse struct {
	FileListMode string `json:"fileListMode"`
	FileList     []struct {
		FileName       string `json:"fileName"`
		TrackType      string `json:"trackType"`
		Uid            string `json:"uid"`
		MixedAllUser   bool   `json:"mixedAllUser"`
		IsPlayable     bool   `json:"isPlayable"`
		SliceStartTime int64  `json:"sliceStartTime"`
	} `json:"fileList"`
	UploadingStatus string `json:"uploadingStatus"`
}

func (s *Update) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateReqBody) (*UpdateResp, error) {
	path := s.buildPath(resourceID, sid, mode)

	responseData, err := s.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}
	var resp UpdateResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse UpdateSuccessResp
		if err = responseData.UnmarshallToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResponse
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

func (s *Update) DoWebRecording(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *UpdateWebRecordingClientRequest) (*UpdateResp, error) {
	return s.Do(ctx, resourceID, sid, WebMode, &UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &UpdateClientRequest{
			WebRecordingConfig: clientRequest.WebRecordingConfig,
			RtmpPublishConfig:  clientRequest.RtmpPublishConfig,
		},
	})
}
