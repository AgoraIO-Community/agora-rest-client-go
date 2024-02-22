package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Stop struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

// BuildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/stop
func (s *Stop) BuildPath(resourceID string, sid string, mode string) string {
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
	SuccessResp StopSuccessResp
}

type StopSuccessResp struct {
	ResourceId string `json:"resourceId"`
	SID        string `json:"sid"`

	serverResponseMode                  StopRespServerResponseMode
	individualRecordingServerResponse   *StopIndividualRecordingServerResponse
	individualVideoScreenshotResponse   *StopIndividualVideoScreenshotServerResponse
	mixRecordingHLSServerResponse       *StopMixRecordingHLSServerResponse
	mixRecordingHLSAndMP4ServerResponse *StopMixRecordingHLSAndMP4ServerResponse
	webRecordingServerResponse          *StopWebRecordingServerResponse
}

type StopWebRecordingResp struct {
	Response
	SuccessResp StopWebRecordingSuccessResp
}

type StopWebRecordingSuccessResp struct {
	ResourceId     string
	SID            string
	ServerResponse StopWebRecordingServerResponse
}

type StopIndividualRecordingServerResponse struct {
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

type StopIndividualVideoScreenshotServerResponse struct {
	UploadingStatus string `json:"uploadingStatus"`
}

type StopMixRecordingHLSServerResponse struct {
	FileListMode    string `json:"fileListMode"`
	FileList        string `json:"fileList"`
	UploadingStatus string `json:"uploadingStatus"`
}

type StopMixRecordingHLSAndMP4ServerResponse struct {
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

type StopWebRecordingServerResponse struct {
	ExtensionServiceState []struct {
		Payload struct {
			UploadingStatus string `json:"uploadingStatus"`
			FileList        []struct {
				FileName       string `json:"fileName"`
				TrackType      string `json:"trackType"`
				Uid            string `json:"uid"`
				MixedAllUser   bool   `json:"mixedAllUser"`
				IsPlayable     bool   `json:"isPlayable"`
				SliceStartTime int64  `json:"sliceStartTime"`
			} `json:"fileList"`
			Onhold bool   `json:"onhold"`
			State  string `json:"state"`
		} `json:"payload"`
		ServiceName string `json:"serviceName"`
		ExitReason  string `json:"exit_reason"`
	} `json:"extensionServiceState"`
}

func (s *StopSuccessResp) GetIndividualRecordingServerResponse() *StopIndividualRecordingServerResponse {
	if s.individualRecordingServerResponse == nil {
		return nil
	}
	return s.individualRecordingServerResponse
}

func (s *StopSuccessResp) GetIndividualVideoScreenshotServerResponse() *StopIndividualVideoScreenshotServerResponse {
	if s.individualVideoScreenshotResponse == nil {
		return nil
	}
	return s.individualVideoScreenshotResponse
}

func (s *StopSuccessResp) GetMixRecordingHLSServerResponse() *StopMixRecordingHLSServerResponse {
	if s.mixRecordingHLSServerResponse == nil {
		return nil
	}
	return s.mixRecordingHLSServerResponse
}

func (s *StopSuccessResp) GetMixRecordingHLSAndMP4ServerResponse() *StopMixRecordingHLSAndMP4ServerResponse {
	if s.mixRecordingHLSAndMP4ServerResponse == nil {
		return nil
	}
	return s.mixRecordingHLSAndMP4ServerResponse
}

func (s *StopSuccessResp) GetWebRecordingServerResponse() *StopWebRecordingServerResponse {
	if s.webRecordingServerResponse == nil {
		return nil
	}
	return s.webRecordingServerResponse
}
func (s *StopSuccessResp) GetServerResponseMode() StopRespServerResponseMode {
	return s.serverResponseMode
}

func (s *StopSuccessResp) setServerResponse(rawBody []byte, mode string) error {
	serverResponseMode := StopServerResponseUnknownMode
	switch mode {
	case IndividualMode:
		fileListMode := gjson.GetBytes(rawBody, "serverResponse.fileListMode")
		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		if fileListMode.Exists() && fileListMode.String() == "json" {
			serverResponseMode = StopIndividualRecordingServerResponseMode

			var resp StopIndividualRecordingServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			s.individualRecordingServerResponse = &resp

		} else {
			serverResponseMode = StopIndividualVideoScreenshotServerResponseMode
			var resp StopIndividualVideoScreenshotServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			s.individualVideoScreenshotResponse = &resp
		}

	case MixMode:
		fileListMode := gjson.GetBytes(rawBody, "serverResponse.fileListMode")
		if !fileListMode.Exists() {
			break
		}

		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		switch fileListMode.String() {
		case "string":
			serverResponseMode = StopMixRecordingHlsServerResponseMode
			var resp StopMixRecordingHLSServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			s.mixRecordingHLSServerResponse = &resp
		case "json":
			serverResponseMode = StopMixRecordingHlsAndMp4ServerResponseMode
			var resp StopMixRecordingHLSAndMP4ServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			s.mixRecordingHLSAndMP4ServerResponse = &resp
		default:
			return errors.New("unknown fileList mode")
		}

	case WebMode:
		serverResponseMode = StopWebRecordingServerResponseMode
		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		var resp StopWebRecordingServerResponse
		if serverResponse.Exists() {
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
		}
		s.webRecordingServerResponse = &resp
	default:
		return errors.New("unknown mode")
	}
	s.serverResponseMode = serverResponseMode
	return nil
}

func (s *Stop) Do(ctx context.Context, resourceID string, sid string, mode string, payload *StopReqBody) (*StopResp, error) {
	path := s.BuildPath(resourceID, sid, mode)

	responseData, err := s.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp StopResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResp StopSuccessResp
		if err = responseData.UnmarshallToTarget(&successResp); err != nil {
			return nil, err
		}
		resp.SuccessResp = successResp
		if err = resp.SuccessResp.setServerResponse(responseData.RawBody, mode); err != nil {
			return nil, err
		}
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

func (s *Stop) DoWebRecording(ctx context.Context, resourceID string, sid string, payload *StopReqBody) (*StopWebRecordingResp, error) {
	mode := WebMode
	resp, err := s.Do(ctx, resourceID, sid, mode, payload)
	if err != nil {
		return nil, err
	}

	var webResp StopWebRecordingResp

	webResp.Response = resp.Response
	if resp.IsSuccess() {
		successResp := resp.SuccessResp
		webResp.SuccessResp = StopWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			SID:            successResp.SID,
			ServerResponse: *successResp.GetWebRecordingServerResponse(),
		}
	}

	return &webResp, nil
}
