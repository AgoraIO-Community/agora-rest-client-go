package individualrecording

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type AcquireIndividualRecodingClientRequest struct {
	// 云端录制 RESTful API 的调用时效。从成功开启云端录制并获得 sid （录制 ID）后开始计算。单位为小时，默认值为 72 小时。
	ResourceExpiredHour int

	// 另一路或几路录制任务的 resourceId。该字段用于排除指定的录制资源，以便新发起的录制任务可以使用新区域的资源，实现跨区域多路录制。
	ExcludeResourceIds []string

	// 指定使用某个区域的资源进行录制。支持取值如下：
	//
	// 0: 根据发起请求的区域就近调用资源。
	//
	// 1: 中国。
	//
	// 2: 东南亚。
	//
	// 3: 欧洲。
	//
	// 4: 北美。
	RegionAffinity int

	// StartParameter 设置该字段后，可以提升可用性并优化负载均衡。
	//
	// 注意：如果填写该字段，则必须确保 startParameter object 和后续 start 请求中填写的 clientRequest object 完全一致，
	// 且取值合法，否则 start 请求会收到报错。
	StartParameter *StartIndividualRecordingClientRequest
}

type StartIndividualRecordingClientRequest struct {
	// Token 用于鉴权的动态密钥（Token）。如果你的项目已启用 App 证书，则务必在该字段中传入你项目的动态密钥
	Token string

	// StorageConfig 第三方云存储的配置项
	StorageConfig *api.StorageConfig

	// RecordingConfig 录制的音视频流配置项
	RecordingConfig *api.RecordingConfig

	// RecordingFileConfig 录制文件的配置项
	RecordingFileConfig *api.RecordingFileConfig

	// SnapshotConfig 视频截图的配置项
	SnapshotConfig *api.SnapshotConfig

	// AppsCollection 应用配置项
	AppsCollection *api.AppsCollection

	// TranscodeOptions 延时转码或延时混音下，生成的录制文件的配置项
	TranscodeOptions *api.TranscodeOptions
}

type QueryIndividualRecordingSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryIndividualRecordingServerResponse
}

type QueryIndividualRecordingResp struct {
	api.Response
	SuccessResponse QueryIndividualRecordingSuccessResp
}

type QueryIndividualRecordingVideoScreenshotSuccessResp struct {
	ResourceId     string
	Sid            string
	ServerResponse *api.QueryIndividualVideoScreenshotServerResponse
}

type QueryIndividualRecordingVideoScreenshotResp struct {
	api.Response
	SuccessResponse QueryIndividualRecordingVideoScreenshotSuccessResp
}

type UpdateIndividualRecordingClientRequest struct {
	StreamSubscribe *api.UpdateStreamSubscribe
}
