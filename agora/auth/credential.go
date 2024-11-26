package auth

import (
	"net/http"
)

type Credential interface {
	Name() string
	SetAuth(r *http.Request)
}
