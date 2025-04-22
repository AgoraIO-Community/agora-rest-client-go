package req

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief Client request for acquiring individual recording resources.
//
// @since v0.8.0
type AcquireIndividualRecordingClientRequest struct {
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
	StartParameter *StartIndividualRecordingClientRequest
}

// @brief Client request for starting individual recording.
//
// @since v0.8.0
type StartIndividualRecordingClientRequest struct {
	// Agora App Token.(Optional)
	Token string

	// Configuration for third-party cloud storage.(Required)
	StorageConfig *api.StorageConfig

	// Configuration for recorded audio and video streams.(Optional)
	RecordingConfig *api.RecordingConfig

	// Configuration for recorded files.(Optional)
	RecordingFileConfig *api.RecordingFileConfig

	// Configurations for screenshot capture.(Optional)
	SnapshotConfig *api.SnapshotConfig
}

// @brief Client request for updating individual recording.
//
// @since v0.8.0
type UpdateIndividualRecordingClientRequest struct {
	// Update subscription lists.(Optional)
	StreamSubscribe *api.UpdateStreamSubscribe
}
