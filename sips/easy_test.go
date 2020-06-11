package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	at := assert.New(t)

	uri := NewSIPRequest("")
	uri.User = "01798300765"
	uri.Host = "87.252.61.202"

	req := MakeRequest(SIPMethodBye, uri, nil, nil)
	at.Equal("BYE sip:01798300765@87.252.61.202 SIP/2.0\r\ncontent-length: 0\r\n", req.String())

	req = MakeRequest(SIPMethodBye, uri, nil, []byte{1, 2, 3})
	at.Equal("BYE sip:01798300765@87.252.61.202 SIP/2.0\r\ncontent-length: 3\r\n\r\n\x01\x02\x03", req.String())

}

func TestMakeResponse(t *testing.T) {
	at := assert.New(t)

	hdr := NewSIPHeader()
	hdr.SetContentLength(1024)
	resp := MakeResponse(StatusOK, hdr, nil)

	at.Equal("SIP/2.0 200 OK\r\ncontent-length: 0\r\n", resp.String())
}
