package core

import (
	"net/http"
)

type CredentialType int

const (
	BasicAuth CredentialType = iota
)

type Credential interface {
	Type() CredentialType
	Visit(r *http.Request)
}

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

func (b *BasicAuthCredential) Type() CredentialType {
	return BasicAuth
}

func (b *BasicAuthCredential) Visit(r *http.Request) {
	r.SetBasicAuth(b.Username, b.Password)
}
