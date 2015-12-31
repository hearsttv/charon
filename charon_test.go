package charon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Charon(t *testing.T) {
	var authState, failRequiredCalled, failProhibitedCalled, handlerCalled bool
	authedFunc := func(r *http.Request) bool { return authState }
	failRequiredFunc := func(rw http.ResponseWriter, r *http.Request) { failRequiredCalled = true }
	failProhibitedFunc := func(rw http.ResponseWriter, r *http.Request) { failProhibitedCalled = true }
	handlerFunc := func(rw http.ResponseWriter, r *http.Request) { handlerCalled = true }

	ch := New(authedFunc, failRequiredFunc, failProhibitedFunc)

	// authed tests
	authedTestFunc := ch.Authenticated(handlerFunc)

	authState = false
	authedTestFunc(httptest.NewRecorder(), &http.Request{})
	assert.True(t, failRequiredCalled)
	assert.False(t, failProhibitedCalled)
	assert.False(t, handlerCalled)

	resetFlags(&authState, &failRequiredCalled, &failProhibitedCalled, &handlerCalled)
	authState = true
	authedTestFunc(httptest.NewRecorder(), &http.Request{})
	assert.False(t, failRequiredCalled)
	assert.False(t, failProhibitedCalled)
	assert.True(t, handlerCalled)

	// unauthed tests
	resetFlags(&authState, &failRequiredCalled, &failProhibitedCalled, &handlerCalled)
	unauthedTestFunc := ch.Unauthenticated(handlerFunc)

	authState = false
	unauthedTestFunc(httptest.NewRecorder(), &http.Request{})
	assert.False(t, failRequiredCalled)
	assert.False(t, failProhibitedCalled)
	assert.True(t, handlerCalled)

	resetFlags(&authState, &failRequiredCalled, &failProhibitedCalled, &handlerCalled)
	authState = true
	unauthedTestFunc(httptest.NewRecorder(), &http.Request{})
	assert.False(t, failRequiredCalled)
	assert.True(t, failProhibitedCalled)
	assert.False(t, handlerCalled)
}

func resetFlags(a, b, c, d *bool) {
	*a, *b, *c, *d = false, false, false, false
}
