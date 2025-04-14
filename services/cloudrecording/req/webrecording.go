package req

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief Client request for acquiring web recording resources.
//
// @since v0.8.0
type AcquireWebRecodingClientRequest struct {
	// The validity period for calling the cloud recording RESTful API.(required)
	//
	// Start calculating after you successfully initiate the cloud recording service and obtain the sid (Recording ID).
	//
	// The calculation unit is hours.
	//
	// The value range is [1,720]. The default value is 72.
	ResourceExpiredHour int

	// The resourceId of another or several other recording tasks.(optional)
	ExcludeResourceIds []string

	// Specify regions that the cloud recording service can access.(optional)
	//
	// The region can be set to:
	//
	//  - 0: Closest to request origin (default)
	// 	- 1: China
	// 	- 2: Southeast Asia
	// 	- 3: Europe
	// 	- 4: North America
	RegionAffinity int

	// StartParameter improves availability and optimizes load balancing.
	StartParameter *StartWebRecordingClientRequest
}

// @brief Client request for starting web page recording.
//
// @since v0.8.0
type StartWebRecordingClientRequest struct {
	// Configuration for recorded files.(Optional)
	RecordingFileConfig *api.RecordingFileConfig

	// Configuration for third-party cloud storage.(Optional)
	StorageConfig *api.StorageConfig

	// Configurations for extended services.(Optional)
	ExtensionServiceConfig *api.ExtensionServiceConfig
}

// @brief Client request for updating web page recording.
//
// @since v0.8.0
type UpdateWebRecordingClientRequest struct {
	// Used to update the web page recording configurations.(Optional)
	WebRecordingConfig *api.UpdateWebRecordingConfig
	// Used to update the configurations for pushing web page recording to the CDN.(Optional)
	RtmpPublishConfig *api.UpdateRtmpPublishConfig
}
