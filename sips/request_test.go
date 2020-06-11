package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_ParseURI(t *testing.T) {
	at := assert.New(t)

	r := NewSIPRequest("sip:01798300765@87.252.61.202:5060;user=phone")

	at.Equal("sip", r.URIType)
	at.Equal("01798300765", r.User)
	at.Equal("87.252.61.202", r.Host)
	at.Equal("5060", r.Port)
	at.Equal("phone", r.Args["user"])

	r = NewSIPRequest("sip:01798300765@87.252.61.202:5060")
	at.Equal("sip", r.URIType)
	at.Equal("01798300765", r.User)
	at.Equal("87.252.61.202", r.Host)
	at.Equal("5060", r.Port)
	at.Len(r.Args, 0)

	r = NewSIPRequest("sip:0179@8300765@87.252.61.202:5060")
	at.Equal("sip", r.URIType)
	at.Equal("0179", r.User)
	at.Equal("8300765@87.252.61.202", r.Host)
	at.Equal("5060", r.Port)
	at.Len(r.Args, 0)

	r = NewSIPRequest("sip:user:password@host:port;uri-parameters?headers")
	at.Equal("sip", r.URIType)
	at.Equal("user:password", r.User)
	at.Equal("host", r.Host)
	at.Equal("port", r.Port)
	at.Equal("", r.Args["uri-parameters?headers"])
}

func TestRequest_String(t *testing.T) {
	at := assert.New(t)

	r := &SIPRequest{
		URIType: "sip",
		User:    "01798300765",
		Host:    "87.252.61.202",
		Port:    "5060",
		Args:    ParseArgs(";user=phone"),
	}
	at.Equal("sip:01798300765@87.252.61.202:5060;user=phone", r.String())

	r = &SIPRequest{
		URIType: "sip",
		User:    "01798300765",
		Host:    "87.252.61.202",
		Port:    "5060",
	}
	at.Equal("sip:01798300765@87.252.61.202:5060", r.String())

	r = &SIPRequest{
		URIType: "sip",
		User:    "01798@300765",
		Host:    "87.252.61.202",
		Port:    "5060",
	}
	at.Equal("sip:01798@300765@87.252.61.202:5060", r.String())
}

func TestSIP_ParseFirstLine(t *testing.T) {
	at := assert.New(t)

	sip := NewSIP()
	err := sip.ParseFirstLine([]byte("INVITE sip:01798300765@87.252.61.202:5060;user=phone SIP/2.0"))
	at.Nil(err)
	at.Equal(SIPMethodInvite, sip.Method)
	at.Equal("sip:01798300765@87.252.61.202:5060;user=phone", string(sip.Request.Src))
	at.NotNil(sip.Request)
	at.False(sip.IsResponse)
	at.Equal(SIPVersion2, sip.Version)

	sip = NewSIP()
	err = sip.ParseFirstLine([]byte("SIP/2.0 200 OK"))
	at.Nil(err)
	at.Equal(SIPMethodNil, sip.Method)
	at.Nil(sip.Request)
	at.True(sip.IsResponse)
	at.Equal(SIPVersion2, sip.Version)
	at.Equal(StatusOK, sip.ResponseCode)
	at.Equal("OK", sip.ResponseCode.Text())
}
