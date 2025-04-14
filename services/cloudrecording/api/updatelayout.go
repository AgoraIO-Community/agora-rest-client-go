package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
)

type UpdateLayout struct {
	client     client.Client
	prefixPath string // /apps/{appid}/cloud_recording
}

func NewUpdateLayout(client client.Client, prefixPath string) *UpdateLayout {
	return &UpdateLayout{client: client, prefixPath: prefixPath}
}

// buildPath returns the request path.
// /v1/apps/{appid}/cloud_recording/resourceid/{resourceid}/sid/{sid}/mode/{mode}/updateLayout
func (u *UpdateLayout) buildPath(resourceID string, sid string, mode string) string {
	return u.prefixPath + "/resourceid/" + resourceID + "/sid/" + sid + "/mode/" + mode + "/updateLayout"
}

type UpdateLayoutReqBody struct {
	Cname         string                     `json:"cname"`
	Uid           string                     `json:"uid"`
	ClientRequest *UpdateLayoutClientRequest `json:"clientRequest"`
}

type UpdateLayoutClientRequest struct {
	MaxResolutionUID           string               `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout           int                  `json:"mixedVideoLayout"`
	BackgroundColor            string               `json:"backgroundColor,omitempty"`
	BackgroundImage            string               `json:"backgroundImage,omitempty"`
	DefaultUserBackgroundImage string               `json:"defaultUserBackgroundImage,omitempty"`
	LayoutConfig               []UpdateLayoutConfig `json:"layoutConfig,omitempty"`
	BackgroundConfig           []BackgroundConfig   `json:"backgroundConfig,omitempty"`
}

// @brief The layout configuration.
//
// @since v0.8.0
type UpdateLayoutConfig struct {
	// The content of the string is the UID of the user to be displayed in the area, 32-bit unsigned integer.(Optional)
	//
	// If the UID is not specified, the screen settings in layoutConfig will be matched automatically in the order that users join the channel.
	UID string `json:"uid"`

	// The relative value of the horizontal coordinate of the upper-left corner of the screen, accurate to six decimal places.(Required)
	//
	// Layout from left to right, with 0.0 at the far left and 1.0 at the far right.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	XAxis float32 `json:"x_axis"`

	// The relative value of the vertical coordinate of the upper-left corner of this screen in the screen, accurate to six decimal places.(Required)
	//
	// Layout from top to bottom, with 0.0 at the top and 1.0 at the bottom.
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	YAxis float32 `json:"y_axis"`

	// The relative value of the width of this screen, accurate to six decimal places.(Required)
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	Width float32 `json:"width"`

	// The relative value of the height of this screen, accurate to six decimal places.(Required)
	//
	// This field can also be set to the integer 0 or 1.
	//
	// The value range is [0,1].
	Height float32 `json:"height"`

	// The transparency of the user's video window.(Optional)
	//
	// Accurate to six decimal places.
	//
	// 0.0 means the user's video window is transparent, and 1.0 indicates that it is completely opaque.
	//
	// The value range is [0,1].
	//
	// The default value is 1.
	Alpha float32 `json:"alpha"`

	// The display mode of users' video windows.(Optional)
	//
	// The value can be set to:
	//
	//  - 0: cropped mode.(Default)
	//     Prioritize to ensure the screen is filled.
	//     The video window size is proportionally scaled until it fills the screen.
	//     If the video's length and width differ from the video window,
	//     the video stream will be cropped from its edges to fit the window, under the aspect ratio set for the video window.
	//  - 1: Fit mode.
	//     Prioritize to ensure that all video content is displayed.
	//     The video size is scaled proportionally until one side of the video window is aligned with the screen border.
	//     If the video scale does not comply with the window size, the video will be scaled to fill the screen while maintaining its aspect ratio.
	//     This scaling may result in a black border around the edges of the video.
	RenderMode int `json:"render_mode"`
}

// @brief Configurations of user's background image.
//
// @since v0.8.0
type BackgroundConfig struct {
	// The string content is the UID.(Required)
	UID string `json:"uid"`

	// The URL of the user's background image.(Required)
	//
	// After setting the background image, if the user stops sending the video stream for more than 3.5 seconds,
	// the screen will switch to the background image.
	// URL supports the HTTP and HTTPS protocols, and the image formats supported are JPG and BMP.
	// The image size must not exceed 6 MB.
	// The settings will only take effect after the recording service successfully downloads the image;
	// if the download fails, the settings will not take effect.
	// Different field settings may overlap each other.
	ImageURL string `json:"image_url"`

	// The display mode of users' video windows.(Optional)
	//
	// The value can be set to:
	//
	//  - 0: cropped mode.(Default)
	//     Prioritize to ensure the screen is filled.
	//     The video window size is proportionally scaled until it fills the screen.
	//     If the video's length and width differ from the video window,
	//     the video stream will be cropped from its edges to fit the window, under the aspect ratio set for the video window.
	//  - 1: Fit mode.
	//     Prioritize to ensure that all video content is displayed.
	//     The video size is scaled proportionally until one side of the video window is aligned with the screen border.
	RenderMode int `json:"render_mode"`
}

// @brief Successful response returned by the cloud recording UpdateLayout API.
//
// @since v0.8.0
type UpdateLayoutSuccessResp struct {
	// Unique identifier of the resource
	ResourceId string `json:"resourceId"`
	// Unique identifier of the recording session
	Sid string `json:"sid"`
}

// @brief Response returned by the cloud recording UpdateLayout API.
//
// @since v0.8.0
type UpdateLayoutResp struct {
	// Response returned by the cloud recording API, see Response for details
	Response
	// Success response, see UpdateLayoutSuccessResp for details
	SuccessResponse UpdateLayoutSuccessResp
}

func (u *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, mode string, payload *UpdateLayoutReqBody) (*UpdateLayoutResp, error) {
	path := u.buildPath(resourceID, sid, mode)

	responseData, err := u.client.DoREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var resp UpdateLayoutResp
	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse UpdateLayoutSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		resp.SuccessResponse = successResponse
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
