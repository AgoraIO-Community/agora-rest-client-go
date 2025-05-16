package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Query struct {
	baseHandler
}

func NewQuery(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Query {
	return &Query{
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
	QueryRtmpPublishServerResponseMode
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
	rtmpPublishServerResponse           *QueryRtmpPublishServerResponse
}

type QueryResp struct {
	Response
	SuccessResponse QuerySuccessResp
}

// @brief Server response returned by the individual recording Query API.
//
// @since v0.8.0
type QueryIndividualRecordingServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`

	// The data format of the fileList field:
	//
	//  - "string": fileList is of String type. In composite recording mode,
	//     if avFileType is set to ["hls"], fileListMode is "string".
	//
	//  - "json": fileList is a JSON Array. When avFileType is set to ["hls","mp4"]
	//     in the individual or composite recording mode, fileListMode is set to "json".
	FileListMode string `json:"fileListMode"`

	// The file list.
	FileList []struct {
		// The file names of the M3U8 and MP4 files generated during recording.
		FileName string `json:"fileName"`

		// The recording file type.
		//
		//  - "audio": Audio-only files.
		//
		//  - "video": Video-only files.
		//
		//  - "audio_and_video": audio and video files
		TrackType string `json:"trackType"`

		// User UID, indicating which user's audio or video stream is being recorded.
		//
		// In composite recording mode, the uid is "0".
		Uid string `json:"uid"`

		// Whether the users were recorded separately.
		//
		//  - true: All users are recorded in a single file.
		//
		//  - false: Each user is recorded separately.
		MixedAllUser bool `json:"mixedAllUser"`

		// Whether or not can be played online.
		//
		//  - true: The file can be played online.
		//
		//  - false: The file cannot be played online.
		IsPlayable bool `json:"isPlayable"`

		// The recording start time of the file, the Unix timestamp, in seconds.
		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	// The recording start time of the file, the Unix timestamp, in seconds.
	SliceStartTime int64 `json:"sliceStartTime"`
}

// @brief Server response returned by the individual recording QueryVideoScreenshot API.
//
// @since v0.8.0
type QueryIndividualVideoScreenshotServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`

	// The recording start time of the file, the Unix timestamp, in seconds.
	SliceStartTime int64 `json:"sliceStartTime"`
}

// @brief Server response returned by the mix recording QueryHLS API.
//
// @since v0.8.0
type QueryMixRecordingHLSServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`

	// The data format of the fileList field:
	//
	//  - "string": fileList is of String type. In composite recording mode,
	//     if avFileType is set to ["hls"], fileListMode is "string".
	//
	//  - "json": fileList is a JSON Array. When avFileType is set to ["hls","mp4"]
	//     in the individual or composite recording mode, fileListMode is set to "json".
	FileListMode string `json:"fileListMode"`

	// The file list.
	FileList string `json:"fileList"`

	// The recording start time of the file, the Unix timestamp, in seconds.
	SliceStartTime int64 `json:"sliceStartTime"`
}

// @brief Server response returned by the mix recording QueryHLSAndMP4 API.
//
// @since v0.8.0
type QueryMixRecordingHLSAndMP4ServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`

	// The data format of the fileList field:
	//
	//  - "string": fileList is of String type. In composite recording mode,
	//     if avFileType is set to ["hls"], fileListMode is "string".
	//
	//  - "json": fileList is a JSON Array. When avFileType is set to ["hls","mp4"]
	//     in the individual or composite recording mode, fileListMode is set to "json".
	FileListMode string `json:"fileListMode"`

	// The file list.
	FileList []struct {
		// The file names of the M3U8 and MP4 files generated during recording.
		FileName string `json:"fileName"`

		// The recording file type.
		//
		//  - "audio": Audio-only files.
		//
		//  - "video": Video-only files.
		//
		//  - "audio_and_video": audio and video files
		TrackType string `json:"trackType"`

		// User UID, indicating which user's audio or video stream is being recorded.
		//
		// In composite recording mode, the uid is "0".
		Uid string `json:"uid"`

		// Whether the users were recorded separately.
		//
		//  - true: All users are recorded in a single file.
		//
		//  - false: Each user is recorded separately.
		MixedAllUser bool `json:"mixedAllUser"`

		// Whether or not can be played online.
		//
		//  - true: The file can be played online.
		//
		//  - false: The file cannot be played online.
		IsPlayable bool `json:"isPlayable"`

		// The recording start time of the file, the Unix timestamp, in seconds.
		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	// The recording start time of the file, the Unix timestamp, in seconds.
	SliceStartTime int64 `json:"sliceStartTime"`
}

// @brief Server response returned by the web recording Query API.
//
// @since v0.8.0
type QueryWebRecordingServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`
	// Extension service state
	ExtensionServiceState []struct {
		// Extension service payload
		Payload struct {
			// File list
			FileList []struct {
				// The file names of the M3U8 and MP4 files generated during recording.
				Filename string `json:"filename"`

				// The recording start time of the file, the Unix timestamp, in seconds.
				SliceStartTime int64 `json:"sliceStartTime"`
			} `json:"fileList"`
			// Whether the page recording is in pause state:
			//
			//  - true: In pause state.
			//
			//  - false: The page recording is running.
			Onhold bool `json:"onhold"`

			// The status of uploading subscription content to the extension service:
			//
			//  - "init": The service is initializing.
			//
			//  - "inProgress": The service has started and is currently in progress.
			//
			//  - "exit": Service exits.
			State string `json:"state"`

			// The status of the push stream to the CDN.
			Outputs []struct {
				// The CDN address to which you push the stream.
				RtmpUrl string `json:"rtmpUrl"`
				// The current status of stream pushing of the web page recording:
				//
				//  - "connecting": Connecting to the CDN server.
				//
				//  - "publishing": The stream pushing is going on.
				//
				//  - "onhold": Set whether to pause the stream pushing.
				//
				//  - "disconnected": Failed to connect to the CDN server. Agora recommends that you change the CDN address to push the stream.
				Status string `json:"status"`
			} `json:"outputs"`
		} `json:"payload"`
		// Name of the extended service:
		//
		//  - "web_recorder_service": Represents the extended service is web page recording.
		//
		//  - "rtmp_publish_service": Represents the extended service is to push web page recording to the CDN.
		ServiceName string `json:"serviceName"`
	} `json:"extensionServiceState"`
}

// @brief Server response returned by the web recording QueryRtmpPublish API.
//
// @since v0.8.0
type QueryRtmpPublishServerResponse struct {
	// Current status of the cloud service:
	//
	//  - 0: Cloud service has not started.
	//
	//  - 1: The cloud service initialization is complete.
	//
	//  - 2: The cloud service components are starting.
	//
	//  - 3: Some cloud service components are ready.
	//
	//  - 4: All cloud service components are ready.
	//
	//  - 5: The cloud service is in progress.
	//
	//  - 6: The cloud service receives the request to stop.
	//
	//  - 7: All components of the cloud service stop.
	//
	//  - 8: The cloud service exits.
	//
	//  - 20: The cloud service exits abnormally.
	Status int `json:"status"`
	// Extension service state
	ExtensionServiceState []struct {
		// Extension service payload
		Payload struct {
			// The status of uploading subscription content to the extension service:
			//
			//  - "init": The service is initializing.
			//
			//  - "inProgress": The service has started and is currently in progress.
			//
			//  - "exit": Service exits.
			State string `json:"state"`

			// The push stream to the CDN output.
			Outputs []struct {
				// The CDN address to which you push the stream.
				RtmpUrl string `json:"rtmpUrl"`
				// The current status of stream pushing of the web page recording:
				//
				//  - "connecting": Connecting to the CDN server.
				//
				//  - "publishing": The stream pushing is going on.
				//
				//  - "onhold": Set whether to pause the stream pushing.
				//
				//  - "disconnected": Failed to connect to the CDN server. Agora recommends that you change the CDN address to push the stream.
				Status string `json:"status"`
			} `json:"outputs"`
		} `json:"payload"`
		// Name of the extended service:
		//
		//  - "web_recorder_service": Represents the extended service is web page recording.
		//
		//  - "rtmp_publish_service": Represents the extended service is to push web page recording to the CDN.
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

func (q *QuerySuccessResp) GetRtmpPublishServiceServerResponse() *QueryRtmpPublishServerResponse {
	return q.rtmpPublishServerResponse
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
		serviceName := gjson.GetBytes(rawBody, "serverResponse.extensionServiceState[*].serviceName")
		serverResponse := gjson.GetBytes(rawBody, "serverResponse")
		switch serviceName.String() {
		case "rtmp_publish_service":
			serverResponseMode = QueryRtmpPublishServerResponseMode
			var resp QueryRtmpPublishServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.rtmpPublishServerResponse = &resp
		case "web_recorder_service":
			serverResponseMode = QueryWebRecordingServerResponseMode
			var resp QueryWebRecordingServerResponse
			if err := json.Unmarshal([]byte(serverResponse.String()), &resp); err != nil {
				return err
			}
			q.webRecordingServerResponse = &resp
		default:
			return errors.New("unknown service name")
		}
	default:
		return errors.New("unknown mode")

	}
	q.serverResponseMode = serverResponseMode
	return nil
}

func (q *Query) Do(ctx context.Context, resourceID string, sid string, mode string) (*QueryResp, error) {
	path := q.buildPath(resourceID, sid, mode)

	responseData, err := doRESTWithRetry(ctx, q.module, q.logger, q.retryCount, q.client, path, http.MethodGet, nil)
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
