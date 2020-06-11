package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSIPAuth(t *testing.T) {
	at := assert.New(t)

	sip := NewSIPAuth(`Digest username="Alice", realm="atlanta.com",
 nonce="84a4cc6f3082121f32b42a2187831a9e",
 response="7587245234b3434cc3412213e5f113a5432"`)

	at.Equal("Alice", sip.Args.Get("username"))
	at.Equal("atlanta.com", sip.Args.Get("realm"))
	at.Equal("84a4cc6f3082121f32b42a2187831a9e", sip.Args.Get("nonce"))
	at.Equal("7587245234b3434cc3412213e5f113a5432", sip.Args.Get("response"))
}
