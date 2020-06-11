package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSIPVia(t *testing.T) {
	at := assert.New(t)

	sv := NewSIPVia("SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK352699498")
	at.Equal("UDP", sv.Trans)
	at.Equal("192.168.1.102", sv.Host)
	at.Equal("5060", sv.Port)
	at.Equal("", sv.Args["rport"])
	at.Equal("z9hG4bK352699498", sv.Args["branch"])
}

func TestSIPVia_Bytes(t *testing.T) {
	at := assert.New(t)

	s := &SIPVia{
		Trans: "UDP",
		Host:  "192.168.1.102",
		Port:  "5060",
	}

	s.Args = make(Args)
	s.Args.Set("branch", "z9hG4bK352699498")
	at.Equal("SIP/2.0/UDP 192.168.1.102:5060;branch=z9hG4bK352699498", s.String())

	s.Args = make(Args)
	s.Args.Set("rport", "")
	at.Equal("SIP/2.0/UDP 192.168.1.102:5060;rport", s.String())
}
