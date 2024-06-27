package v1

import (
	"net/http"

	client "github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type ErrResponse struct {
	ErrorCode int    `json:"code"` // ErrorCode reference to https://doc.shengwang.cn/api-ref/cloud-recording/restful/response-code
	Reason    string `json:"reason"`
}

type Response struct {
	*client.BaseResponse
	ErrResponse ErrResponse
}

func (b Response) IsSuccess() bool {
	if b.BaseResponse != nil {
		return b.HttpStatusCode == http.StatusOK
	} else {
		return false
	}
}
