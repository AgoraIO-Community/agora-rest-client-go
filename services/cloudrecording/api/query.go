package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type Query struct {
	client     client.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

func NewQuery(client client.Client, prefixPath string) *Query {
	return &Query{client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/query
func (q *Query) buildPath(resourceID string, sid string, mode string) string {
	return q.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/query"
}

type QueryRespServerResponseMode int

const (
	QueryServerResponseUnknownMode QueryRespServerResponseMode = iota
	QueryIndividualRecordingServerResponseMode
	QueryIndividualVideoScreenshotServerResponseMode
	QueryMixRecordingHlsServerResponseMode
	QueryMixRecordingHlsAndMp4ServerResponseMode
	QueryWebRecordingServerResponseMode
)

type QuerySuccessResp struct {
	ResourceId string `json:"resourceId"`
	Sid        string `json:"sid"`

	serverResponseMode                  QueryRespServerResponseMode
	individualRecordingServerResponse   *QueryIndividualRecordingServerResponse
	individualVideoScreenshotResponse   *QueryIndividualVideoScreenshotServerResponse
	mixRecordingHLSServerResponse       *QueryMixRecordingHLSServerResponse
	mixRecordingHLSAndMP4ServerResponse *QueryMixRecordingHLSAndMP4ServerResponse
	webRecordingServerResponse          *QueryWebRecordingServerResponse
}

type QueryResp struct {
	Response
	SuccessResponse QuerySuccessResp
}

type QueryIndividualRecordingServerResponse struct {
	Status int `json:"status"`

	FileListMode string `json:"fileListMode"`

	FileList []struct {
		FileName string `json:"fileName"`

		TrackType string `json:"trackType"`

		Uid          string `json:"uid"`
		MixedAllUser bool   `json:"mixedAllUser"`

		IsPlayable bool `json:"isPlayable"`

		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryIndividualVideoScreenshotServerResponse struct {
	Status int `json:"status"`

	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryMixRecordingHLSServerResponse struct {
	Status int `json:"status"`

	FileListMode string `json:"fileListMode"`

	FileList string `json:"fileList"`

	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryMixRecordingHLSAndMP4ServerResponse struct {
	Status int `json:"status"`

	FileListMode string `json:"fileListMode"`

	FileList []struct {
		FileName string `json:"fileName"`

		TrackType string `json:"trackType"`

		Uid string `json:"uid"`

		MixedAllUser bool `json:"mixedAllUser"`

		IsPlayable bool `json:"isPlayable"`

		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryWebRecordingServerResponse struct {
	Status                int `json:"status"`
	ExtensionServiceState []struct {
		Payload struct {
			FileList []struct {
				Filename string `json:"filename"`

				SliceStartTime int64 `json:"sliceStartTime"`
			} `json:"fileList"`

			Onhold bool `json:"onhold"`

			State string `json:"state"`

			Outputs []struct {
				RtmpUrl string `json:"rtmpUrl"`


			} `json:"outputs"`
		} `json:"payload"`
		ServiceName string `json:"serviceName"`
	} `json:"extensionServiceState"`
}

func (q *QuerySuccessResp) GetIndividualRecordingServerResponse() *QueryIndividualRecordingServerResponse {
	return q.individualRecordingServerResponse
}

func (q *QuerySuccessResp) GetIndividualVideoScreenshotServerResponse() *QueryIndividualVideoScreenshotServerResponse {
	return q.individualVideoScreenshotResponse
}

func (q *QuerySuccessResp) GetMixRecordingHLSServerResponse() *QueryMixRecordingHLSServerResponse {
	return q.mixRecordingHLSServerResponse
}

func (q *QuerySuccessResp) GetMixRecordingHLSAndMP4ServerResponse() *QueryMixRecordingHLSAndMP4ServerResponse {
	return q.mixRecordingHLSAndMP4ServerResponse
}

func (q *QuerySuccessResp) GetWebRecording2CDNServerResponse() *QueryWebRecordingServerResponse {
	return q.webRecordingServerResponse
}

func (q *QuerySuccessResp) GetServerResponseMode() QueryRespServerResponseMode {
	return q.serverResponseMode
}

func (q *QuerySuccessResp) setServerResponse(rawBody []byte, mode string) error {
	serverResponseMode := QueryServerResponseUnknownMode
	switch mode {
	case IndividualMode:
		fileListMode := gjson.GetBytes(rawBody, "serverResponse.fileListMode")
		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		if fileListMode.Exists() && fileListMode.String() == "json" {
			serverResponseMode = QueryIndividualRecordingServerResponseMode
			var resp QueryIndividualRecordingServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.individualRecordingServerResponse = &resp

		} else {
			serverResponseMode = QueryIndividualVideoScreenshotServerResponseMode
			var resp QueryIndividualVideoScreenshotServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.individualVideoScreenshotResponse = &resp
		}

	case MixMode:
		fileListMode := gjson.GetBytes(rawBody, "serverResponse.fileListMode")
		if !fileListMode.Exists() {
			break
		}

		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		switch fileListMode.String() {
		case "string":
			serverResponseMode = QueryMixRecordingHlsServerResponseMode
			var resp QueryMixRecordingHLSServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.mixRecordingHLSServerResponse = &resp
		case "json":
			serverResponseMode = QueryMixRecordingHlsAndMp4ServerResponseMode
			var resp QueryMixRecordingHLSAndMP4ServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.mixRecordingHLSAndMP4ServerResponse = &resp
		default:
			return errors.New("unknown fileList mode")
		}

	case WebMode:
		serverResponseMode = QueryWebRecordingServerResponseMode
		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		var resp QueryWebRecordingServerResponse
		if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
			return err
		}
		q.webRecordingServerResponse = &resp
	default:
		return errors.New("unknown mode")

	}
	q.serverResponseMode = serverResponseMode
	return nil
}

func (q *Query) Do(ctx context.Context, resourceID string, sid string, mode string) (*QueryResp, error) {
	path := q.buildPath(resourceID, sid, mode)

	responseData, err := q.client.DoREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}
	var resp QueryResp

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse QuerySuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResponse = successResponse
		if err = resp.SuccessResponse.setServerResponse(responseData.RawBody, mode); err != nil {
			return nil, err
		}
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
