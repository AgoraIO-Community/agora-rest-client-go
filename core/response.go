package core

import (
	"encoding/json"
)

type ResponseInterface interface {
	IsSuccess() bool
}

type BaseResponse struct {
	RawBody        []byte
	HttpStatusCode int
}

// UnmarshallToTarget unmarshall body into target var
// successful if err is nil
func (r *BaseResponse) UnmarshallToTarget(target interface{}) error {
	err := json.Unmarshal(r.RawBody, target)
	if err != nil {
		return err
	}
	return nil
}
