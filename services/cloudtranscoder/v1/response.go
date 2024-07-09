package v1

import (
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type ErrResponse struct {
	Message string `json:"message"`
}

type Response struct {
	*core.BaseResponse
	ErrResponse ErrResponse
}

func (b Response) IsSuccess() bool {
	if b.BaseResponse != nil {
		return b.HttpStatusCode == http.StatusOK
	} else {
		return false
	}
}

