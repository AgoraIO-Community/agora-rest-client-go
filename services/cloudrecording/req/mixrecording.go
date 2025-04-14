package req

import "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"

// @brief Client request for acquiring mix recording resources.
//
// @since v0.8.0
type AcquireMixRecodingClientRequest struct {
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
	StartParameter *StartMixRecordingClientRequest
}

// @brief Client request for starting mix recording.
//
// @since v0.8.0
type StartMixRecordingClientRequest struct {
	// Agora App Token.(Optional)
	Token string

	// Configuration for recorded audio and video streams.(Optional)
	RecordingConfig *api.RecordingConfig

	// Configuration for recorded files.(Optional)
	RecordingFileConfig *api.RecordingFileConfig

	// Configuration for third-party cloud storage.(Optional)
	StorageConfig *api.StorageConfig

	// Configuration for screenshot capture.(Optional)
	SnapshotConfig *api.SnapshotConfig
}

// @brief Client request for updating mix recording.
//
// @since v0.8.0
type UpdateMixRecordingClientRequest struct {
	// Update subscription lists.(Optional)
	StreamSubscribe *api.UpdateStreamSubscribe
}

// @brief Client request for updating the layout of mix recording.
//
// @since v0.8.0
type UpdateLayoutUpdateMixRecordingClientRequest struct {
	// Only need to set it in vertical layout.(Optional)
	//
	// Specify the user ID of the large video window.
	MaxResolutionUID string

	// Composite video layout.(Optional)
	//
	// The value can be set to:
	//
	//  - 0: (Default) Floating layout.
	//     The first user to join the channel will be displayed as a large window, filling the entire canvas.
	//     The video windows of other users will be displayed as small windows, arranged horizontally from bottom to top, up to 4 rows, each with 4 windows.
	//     It supports up to a total of 17 windows of different users' videos.
	//  - 1: Adaptive layout.
	//     Automatically adjust the size of each user's video window according to the number of users, each user's video window size is consistent, and supports up to 17 windows.
	//  - 2: Vertical layout.
	//     The maxResolutionUid is specified to display the large video window on the left side of the screen, and the small video windows of other users are vertically arranged on the right side, with a maximum of two columns, 8 windows per column, supporting up to 17 windows.
	//  - 3: Customized layout.
	//     Set the layoutConfig field to customize the mixed layout.
	MixedVideoLayout int

	// The background color of the video canvas.(Optional)
	//
	// The RGB color table is supported, with strings formatted as a # sign and 6 hexadecimal digits.
	//
	// The default value is "#000000", representing the black color.
	BackgroundColor string

	// The background image of the video canvas.(Optional)
	//
	// The display mode of the background image is set to cropped mode.
	//
	// Cropped mode: Will prioritize to ensure that the screen is filled.
	// The background image size is scaled in equal proportion until the entire screen is filled with the background image.
	// If the length and width of the background image differ from the video window,
	// the background image will be peripherally cropped to fill the window.
	BackgroundImage string

	// The URL of the default user screen background image.(Optional)
	//
	// After configuring this field, when any user stops sending the video streams for more than 3.5 seconds,
	// the screen will switch to the background image;
	// this setting will be overridden if the background image is set separately for a UID.
	DefaultUserBackgroundImage string

	// The mixed video layout of users.(Optional)
	//
	// An array of screen layout settings for each user, supporting up to 17 users.
	LayoutConfig []api.UpdateLayoutConfig

	// The background configuration.(Optional)
	//
	// The backgroundConfig field is used to set the background of the video windows.
	BackgroundConfig []api.BackgroundConfig
}
