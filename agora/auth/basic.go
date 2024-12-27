package auth

import (
	"net/http"
)

type BasicAuthCredential struct {
	username string
	password string
}

// Ensure BasicAuthCredential implements Credential
var _ Credential = (*BasicAuthCredential)(nil)

func NewBasicAuthCredential(username string, password string) *BasicAuthCredential {
	return &BasicAuthCredential{
		username: username,
		password: password,
	}
}

func (b *BasicAuthCredential) Name() string {
	return "basic_auth"
}

func (b *BasicAuthCredential) SetAuth(r *http.Request) {
	r.SetBasicAuth(b.username, b.password)
}
