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
type StopClientRequest struct {
	AsyncStop bool `json:"async_stop"`
}

// @brief StopResp returned by the various of cloud recording scenarios Stop API.
//
// @since v0.8.0
type StopResp struct {
	// Response returned by the cloud recording API, see Response for details
	Response
	// Successful response, see StopSuccessResp for details
	SuccessResponse StopSuccessResp
}

// @brief Successful response returned by the various of cloud recording scenarios Stop API.
//
// @since v0.8.0
type StopSuccessResp struct {
	// Name of the channel to be recorded
	Cname string `json:"cname"`
	// User ID used by the cloud recording service in the RTC channel to identify the recording service in the channel
	UID string `json:"uid"`
	// Unique identifier of the resource
	ResourceId string `json:"resourceId"`
	// Unique identifier of the recording session
	Sid string `json:"sid"`
}

func (s *Stop) Do(ctx context.Context, resourceId string, sid string, mode string, payload *StopReqBody) (*StopResp, error) {
	path := s.buildPath(resourceId, sid, mode)

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
