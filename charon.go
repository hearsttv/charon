package charon

import (
	"net/http"
)

// Charon exposes two wrapper methods for HTTP(S) handlers that enforce auth state one way or the other.
type Charon interface {
	Authenticated(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	Unauthenticated(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	HasPriv(handler func(http.ResponseWriter, *http.Request), privileges ...string) func(http.ResponseWriter, *http.Request)
}

type middleware struct {
	isAuthed           func(r *http.Request) bool
	hasPriv            func(r *http.Request, privileges ...string) bool
	failAuthRequired   func(rw http.ResponseWriter, r *http.Request)
	failNotAuthorized  func(rw http.ResponseWriter, r *http.Request)
	failAuthProhibited func(rw http.ResponseWriter, r *http.Request)
}

// New creates a new Charon object with the specified callback functions
func New(isAuthed func(r *http.Request) bool, failAuthRequired func(rw http.ResponseWriter, r *http.Request), failAuthProhibited func(rw http.ResponseWriter, r *http.Request)) Charon {
	return &middleware{
		isAuthed:           isAuthed,
		failAuthRequired:   failAuthRequired,
		failAuthProhibited: failAuthProhibited,
	}
}

func NewWithRoleSupport(isAuthed func(r *http.Request) bool, hasPriv func(r *http.Request, privileges ...string) bool, failAuthRequired func(rw http.ResponseWriter, r *http.Request), failNotAuthorized func(rw http.ResponseWriter, r *http.Request), failAuthProhibited func(rw http.ResponseWriter, r *http.Request)) Charon {
	return &middleware{
		isAuthed:           isAuthed,
		hasPriv:            hasPriv,
		failAuthRequired:   failAuthRequired,
		failNotAuthorized:  failNotAuthorized,
		failAuthProhibited: failAuthProhibited,
	}
}

// Authenticated wraps a handler, enforcing that it will only be called from an authenticated context
func (m *middleware) Authenticated(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if m.isAuthed(r) {
			handler(rw, r)
		} else {
			m.failAuthRequired(rw, r)
		}
	}
}

// Unauthenticated wraps a handler, enforcing that it will only be called from an unauthenticated context
func (m *middleware) Unauthenticated(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !m.isAuthed(r) {
			handler(rw, r)
		} else {
			m.failAuthProhibited(rw, r)
		}
	}
}

func (m *middleware) HasPriv(handler func(http.ResponseWriter, *http.Request), privileges ...string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !m.isAuthed(r) {
			m.failAuthRequired(rw, r)
		} else if !m.hasPriv(r, privileges...) {
			m.failNotAuthorized(rw, r)
		} else {
			handler(rw, r)
		}
	}
}
