package agora

import (
	"encoding/json"
	"net/http"
)

type ResponseInterface interface {
	IsSuccess() bool
}

type BaseResponse struct {
	RawResponse    *http.Response
	RawBody        []byte
	HttpStatusCode int
}

// UnmarshalToTarget unmarshal body into target var
// successful if err is nil
func (r *BaseResponse) UnmarshalToTarget(target interface{}) error {
	err := json.Unmarshal(r.RawBody, target)
	if err != nil {
		return err
	}
	return nil
}

func (r *BaseResponse) GetRequestID() string {
	if r.RawResponse != nil {
		return r.RawResponse.Header.Get("X-Request-Id")
	} else {
		return ""
	}
}
