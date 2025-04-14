package api

import (
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
)

// @brief Error response returned by the cloud recording API.
//
// @since v0.8.0
type ErrResponse struct {
	// Error code
	ErrorCode int `json:"code"`
	// Reason for the error
	Reason string `json:"reason"`
}

// @brief Response returned by the cloud recording API.
//
// @since v0.8.0
type Response struct {
	// HTTP base response, see agora.BaseResponse for details
	*agora.BaseResponse
	// Error response, see ErrResponse for details
	ErrResponse ErrResponse
}

// @brief Determines whether the response returned by the cloud recording API is successful.
//
// @note If the response is successful, continue to read the data in the successful response; otherwise, read the data in the error response
//
// @return Returns true if successful, otherwise false
//
// @since v0.8.0
func (b Response) IsSuccess() bool {
	if b.BaseResponse != nil {
		return b.HttpStatusCode == http.StatusOK
	} else {
		return false
	}
}
