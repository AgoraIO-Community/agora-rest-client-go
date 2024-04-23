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

type StopIndividualRecordingServerResponse struct {
	// FileListMode fileList 字段的数据格式：
	// "string"：fileList 为 String 类型。合流录制模式下，如果 avFileType 设置为 ["hls"]，fileListMode 为 "string"。
	//
	// "json"：fileList 为 JSON Array 类型。单流或合流录制模式下 avFileType 设置为 ["hls","mp4"] 时，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode"`

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

	// UploadingStatus 当前录制上传的状态：
	//
	// "uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
	//
	// "backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了声网备份云。声网服务器会自动将这部分文件继续上传至指定的第三方云存储。
	//
	// "unknown"：未知状态。
	UploadingStatus string `json:"uploadingStatus"`
}

type StopIndividualVideoScreenshotServerResponse struct {
	// UploadingStatus 当前录制上传的状态：
	//
	// "uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
	//
	// "backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了声网备份云。声网服务器会自动将这部分文件继续上传至指定的第三方云存储。
	//
	// "unknown"：未知状态。
	UploadingStatus string `json:"uploadingStatus"`
}

type StopMixRecordingHLSServerResponse struct {
	// FileListMode fileList 字段的数据格式：
	// "string"：fileList 为 String 类型。合流录制模式下，如果 avFileType 设置为 ["hls"]，fileListMode 为 "string"。
	//
	// "json"：fileList 为 JSON Array 类型。单流或合流录制模式下 avFileType 设置为 ["hls","mp4"] 时，fileListMode 为 "json"。
	FileListMode string `json:"fileListMode"`

	// FileList 录制产生的 M3U8 文件的文件名。
	FileList string `json:"fileList"`

	// UploadingStatus 当前录制上传的状态：
	//
	// "uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
	//
	// "backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了声网备份云。声网服务器会自动将这部分文件继续上传至指定的第三方云存储。
	//
	// "unknown"：未知状态。
	UploadingStatus string `json:"uploadingStatus"`
}

type StopMixRecordingHLSAndMP4ServerResponse struct {
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

	// UploadingStatus 当前录制上传的状态：
	//
	// "uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
	//
	// "backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了声网备份云。声网服务器会自动将这部分文件继续上传至指定的第三方云存储。
	//
	// "unknown"：未知状态。
	UploadingStatus string `json:"uploadingStatus"`
}

type StopWebRecordingServerResponse struct {
	ExtensionServiceState []struct {
		Payload struct {
			// UploadingStatus 当前录制上传的状态：
			//
			// "uploaded"：本次录制的文件已经全部上传至指定的第三方云存储。
			//
			// "backuped"：本次录制的文件已经全部上传完成，但是至少有一个 TS 文件上传到了声网备份云。声网服务器会自动将这部分文件继续上传至指定的第三方云存储。
			//
			// "unknown"：未知状态。
			UploadingStatus string `json:"uploadingStatus"`

			// FileListMode 文件列表
			FileList []struct {
				// FileName 录制产生的 M3U8 文件和 MP4 文件的文件名。
				FileName string `json:"fileName"`

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
		} `json:"payload"`

		// ServiceName 服务类型
		// 服务类型：
		// "upload_service"：上传服务。
		//
		// "web_recorder_service"：页面录制服务。
		ServiceName string `json:"serviceName"`
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
