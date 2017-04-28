# Charon
Simple authed/unauthed request enforcer in Go

Charon is the mythological boatman that offered passage across the river Styx into the underworld. For a price. 

In our world, Charon is a simple mechanism for marking `net/http` hander funcs as only being callable from an authenticated (or unauthenticated) context. This keeps the logic for auth context from being repeated throughout your handlers.

## Installation 

`go get github.com/hearsttv/charon`

`import "github.com/hearsttv/charon"`
