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
	// Status 当前云服务的状态：
	//
	// 0：没有开始云服务。
	//
	// 1：云服务初始化完成。
	//
	// 2：云服务组件开始启动。
	//
	// 3：云服务部分组件启动完成。
	//
	// 4：云服务所有组件启动完成。
	//
	// 5：云服务正在进行中。
	//
	// 6：云服务收到停止请求。
	//
	// 7：云服务所有组件均停止。
	//
	// 8：云服务已退出。
	//
	// 20：云服务异常退出。
	Status int `json:"status"`

	//  FileListMode fileList 字段的数据格式：
	// "string"：fileList 为 String 类型。合流录制模式下，如果 avFileType 设置为 ["hls"]，fileListMode 为 "string"。
	//
	// "json"：fileList 为 JSON Array 类型。单流或合流录制模式下 avFileType 设置为 ["hls","mp4"] 时，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode"`

	// FileList 录制文件列表。
	FileList []struct {
		// FileName 录制产生的 M3U8 文件和 MP4 文件的文件名。
		FileName string `json:"fileName"`

		// TrackType 录制文件的类型:
		//
		// "audio"：纯音频文件。
		//
		// "video"：纯视频文件。
		//
		// "audio_and_video"：音视频文件。
		TrackType string `json:"trackType"`

		// Uid 用户 UID，表示录制的是哪个用户的音频流或视频流。
		//
		// 合流录制模式下，uid 为 "0"。
		Uid string `json:"uid"`

		// MixedAllUser 用户是否是分开录制
		//
		// true：所有用户合并在一个录制文件中。
		//
		// false：每个用户分开录制。
		MixedAllUser bool `json:"mixedAllUser"`

		// IsPlayable 是否可以在线播放。
		//
		// true：可以在线播放。
		//
		// false：无法在线播放。
		IsPlayable bool `json:"isPlayable"`

		// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryIndividualVideoScreenshotServerResponse struct {
	// Status 当前云服务的状态：
	//
	// 0：没有开始云服务。
	//
	// 1：云服务初始化完成。
	//
	// 2：云服务组件开始启动。
	//
	// 3：云服务部分组件启动完成。
	//
	// 4：云服务所有组件启动完成。
	//
	// 5：云服务正在进行中。
	//
	// 6：云服务收到停止请求。
	//
	// 7：云服务所有组件均停止。
	//
	// 8：云服务已退出。
	//
	// 20：云服务异常退出。
	Status int `json:"status"`

	// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryMixRecordingHLSServerResponse struct {
	// Status 当前云服务的状态：
	//
	// 0：没有开始云服务。
	//
	// 1：云服务初始化完成。
	//
	// 2：云服务组件开始启动。
	//
	// 3：云服务部分组件启动完成。
	//
	// 4：云服务所有组件启动完成。
	//
	// 5：云服务正在进行中。
	//
	// 6：云服务收到停止请求。
	//
	// 7：云服务所有组件均停止。
	//
	// 8：云服务已退出。
	//
	// 20：云服务异常退出。
	Status int `json:"status"`

	//  FileListMode fileList 字段的数据格式：
	// "string"：fileList 为 String 类型。合流录制模式下，如果 avFileType 设置为 ["hls"]，fileListMode 为 "string"。
	//
	// "json"：fileList 为 JSON Array 类型。单流或合流录制模式下 avFileType 设置为 ["hls","mp4"] 时，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode"`

	// FileList 录制产生的 M3U8 文件的文件名。
	FileList string `json:"fileList"`

	// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryMixRecordingHLSAndMP4ServerResponse struct {
	// Status 当前云服务的状态：
	//
	// 0：没有开始云服务。
	//
	// 1：云服务初始化完成。
	//
	// 2：云服务组件开始启动。
	//
	// 3：云服务部分组件启动完成。
	//
	// 4：云服务所有组件启动完成。
	//
	// 5：云服务正在进行中。
	//
	// 6：云服务收到停止请求。
	//
	// 7：云服务所有组件均停止。
	//
	// 8：云服务已退出。
	//
	// 20：云服务异常退出。
	Status int `json:"status"`

	// FileListMode fileList 字段的数据格式：
	// "string"：fileList 为 String 类型。合流录制模式下，如果 avFileType 设置为 ["hls"]，fileListMode 为 "string"。
	//
	// "json"：fileList 为 JSON Array 类型。单流或合流录制模式下 avFileType 设置为 ["hls","mp4"] 时，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode"`

	// FileList 录制文件列表。
	FileList []struct {
		// FileName 录制产生的 M3U8 文件和 MP4 文件的文件名。
		FileName string `json:"fileName"`

		// TrackType 录制文件的类型:
		//
		// "audio"：纯音频文件。
		//
		// "video"：纯视频文件。
		//
		// "audio_and_video"：音视频文件。
		TrackType string `json:"trackType"`

		// Uid 用户 UID，表示录制的是哪个用户的音频流或视频流。
		//
		// 合流录制模式下，uid 为 "0"。
		Uid string `json:"uid"`

		// MixedAllUser 用户是否是分开录制
		//
		// true：所有用户合并在一个录制文件中。
		//
		// false：每个用户分开录制。
		MixedAllUser bool `json:"mixedAllUser"`

		// IsPlayable 是否可以在线播放。
		//
		// true：可以在线播放。
		//
		// false：无法在线播放。
		IsPlayable bool `json:"isPlayable"`

		// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
		SliceStartTime int64 `json:"sliceStartTime"`
	} `json:"fileList"`

	// SliceStartTime 录制开始的时间，Unix 时间戳，单位为毫秒。
	SliceStartTime int64 `json:"sliceStartTime"`
}

type QueryWebRecordingServerResponse struct {
	// Status 当前云服务的状态：
	//
	// 0：没有开始云服务。
	//
	// 1：云服务初始化完成。
	//
	// 2：云服务组件开始启动。
	//
	// 3：云服务部分组件启动完成。
	//
	// 4：云服务所有组件启动完成。
	//
	// 5：云服务正在进行中。
	//
	// 6：云服务收到停止请求。
	//
	// 7：云服务所有组件均停止。
	//
	// 8：云服务已退出。
	//
	// 20：云服务异常退出。
	Status                int `json:"status"`
	ExtensionServiceState []struct {
		Payload struct {
			// FileListMode 文件列表
			FileList []struct {
				// FileName 录制产生的 M3U8 文件和 MP4 文件的文件名。
				Filename string `json:"filename"`

				// SliceStartTime 该文件的录制开始时间，Unix 时间戳，单位为毫秒。
				SliceStartTime int64 `json:"sliceStartTime"`
			} `json:"fileList"`

			// Onhold 页面录制是否处于暂停状态：
			//
			// true：处于暂停状态。
			//
			// false：处于运行状态。
			Onhold bool `json:"onhold"`

			// State 将订阅内容上传至扩展服务的状态：
			//
			// "init"：服务正在初始化。
			//
			// "inProgress"：服务启动完成，正在进行中。
			//
			// "exit"：服务退出。
			State string `json:"state"`

			Outputs []struct {
				// RtmpUrl CDN 推流地址。
				RtmpUrl string `json:"rtmpUrl"`

				// Status 页面录制当前的推流状态：
				//
				// "connecting"：正在连接 CDN 服务器。
				//
				// "publishing"：正在推流。
				//
				// "onhold"：设置是否暂停推流。
				//
				// "disconnected"：连接 CDN 服务器失败，声网建议你更换 CDN 推流地址。
				Status string `json:"status"`
			} `json:"outputs"`
		} `json:"payload"`
		// ServiceName 扩展服务的名称：
		//
		// web_recorder_service：代表扩展服务为页面录制。
		//
		// rtmp_publish_service：代表扩展服务为转推页面录制到 CDN。
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
