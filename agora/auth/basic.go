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

// @brief Create a new BasicAuthCredential instance
//
// @note Customer ID and customer secret obtained from the Agora console, used for HTTP basic authentication
//
// @param username Customer ID obtained from the Agora console
//
// @param password Customer secret obtained from the Agora console
//
// @return Returns the BasicAuthCredential instance
//
// @since v0.7.0
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
