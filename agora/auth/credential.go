package auth

import (
	"net/http"
)

type Credential interface {
	Name() string
	Visit(r *http.Request)
}
