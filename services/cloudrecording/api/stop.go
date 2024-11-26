package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type Stop struct {
	client     client.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

func NewStop(client client.Client, prefixPath string) *Stop {
	return &Stop{client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/stop
func (s *Stop) buildPath(resourceID string, sid string, mode string) string {
	return s.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/stop"
}

type StopRespServerResponseMode int

const (
	StopServerResponseUnknownMode StopRespServerResponseMode = iota
	StopIndividualRecordingServerResponseMode
	StopIndividualVideoScreenshotServerResponseMode
	StopMixRecordingHlsServerResponseMode
	StopMixRecordingHlsAndMp4ServerResponseMode
	StopWebRecordingServerResponseMode
)

type StopReqBody struct {
	Cname         string             `json:"cname"`
	Uid           string             `json:"uid"`
	ClientRequest *StopClientRequest `json:"clientRequest"`
}

// StopClientRequest is the request body of stop.
// 当 async_stop 为 true 时，表示异步停止录制。默认值为 false，异步情况下可能会获取不到对应的serverResponse内容
type StopClientRequest struct {
	AsyncStop bool `json:"async_stop"`
}

type StopResp struct {
	Response
	SuccessResponse StopSuccessResp
}

type StopSuccessResp struct {
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

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, mode string, payload *StopReqBody) (*StopResp, error) {
	path := s.buildPath(resourceID, sid, mode)

	responseData, err := s.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp StopResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResp StopSuccessResp
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
