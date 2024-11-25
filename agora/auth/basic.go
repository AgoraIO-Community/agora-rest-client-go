package auth

import (
	"net/http"
)

type BasicAuthCredential struct {
	Username string
	Password string
}

// Ensure BasicAuthCredential implements Credential
var _ Credential = (*BasicAuthCredential)(nil)

func NewBasicAuthCredential(username string, password string) *BasicAuthCredential {
	return &BasicAuthCredential{
		Username: username,
		Password: password,
	}
}

func (b *BasicAuthCredential) Name() string {
	return "basic_auth"
}

func (b *BasicAuthCredential) Visit(r *http.Request) {
	r.SetBasicAuth(b.Username, b.Password)
}
