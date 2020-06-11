package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSIPCseq(t *testing.T) {
	at := assert.New(t)

	sc := NewSIPCseq("4711 INVITE")
	at.Equal("4711", string(sc.ID))
	at.Equal("INVITE", string(sc.Method))
}

func TestSIPCseq_Bytes(t *testing.T) {
	at := assert.New(t)

	sip := NewSIPCseq("")
	sip.ID = "4711"
	sip.Method = "INVITE"

	at.Equal("4711 INVITE", sip.String())
}
