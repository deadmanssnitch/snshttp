package snshttp

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// Option gives a way to customize the SNS handler.
type Option interface {
	apply(*handler)
}

type authOption struct {
	credentials [sha256.Size]byte
}

// noopOption is used when the config for an Option (like authentication) mean
// the option shouldn't be applied.
type noopOption struct{}

func (opt *noopOption) apply(_ *handler) {}

// WithAuthentication protects the webhook endpoint behind basic authentication
// and should be used with HTTPS endpoints as the credentials are basically
// transmitted in plain text. Either the username, password, or both can be
// left blank empty. When both are empty authentication will be disabled.
func WithAuthentication(username string, password string) Option {
	// When neither the username or password is given then we should avoid
	// turning on authentication.
	if username == "" && password == "" {
		return &noopOption{}
	}

	return &authOption{
		// At initialization we generate a SHA256 hash that will be used to check
		// the username and password on each request.
		credentials: sha256.Sum256([]byte(username + ":" + password)),
	}
}

func (opt *authOption) apply(handler *handler) {
	handler.credentials = opt
}

// Check veriies the
func (opt *authOption) Check(req *http.Request) bool {
	// If the option is nil then authentication is not required
	if opt == nil {
		return true
	}

	username, password, provided := req.BasicAuth()
	if !provided {
		// No credentials were provided so fail authentication
		return false
	}

	given := sha256.Sum256([]byte(username + ":" + password))

	// Use a constant time function to avoid timing attacks. This requires both
	// arguments to be the same length (they're both sha256.Size long).
	return subtle.ConstantTimeCompare(opt.credentials[:], given[:]) == 1
}
