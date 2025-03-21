package agora

import (
	"encoding/json"
	"net/http"
)

type ResponseInterface interface {
	IsSuccess() bool
}

// @brief HTTP response
//
// @since v0.7.0
type BaseResponse struct {
	// HTTP Raw response
	RawResponse *http.Response
	// Raw body of the response
	RawBody []byte
	// HTTP status code
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

// @brief Get the request ID from the response
//
// @since v0.7.0
func (r *BaseResponse) GetRequestID() string {
	if r.RawResponse != nil {
		return r.RawResponse.Header.Get("X-Request-Id")
	} else {
		return ""
	}
}
