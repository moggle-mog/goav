package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSIP_Parse(t *testing.T) {
	at := assert.New(t)

	raw := []byte(`MESSAGE sip:34020000002000000001@10.104.157.255:5060 SIP/2.0
Via: SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK1135200957
From: <sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856
To: <sip:34020000002000000001@10.104.157.255:5060>
Call-ID: 1808031968@192.168.1.102
CSeq: 20 MESSAGE
Max-Forwards: 70
User-Agent: SIP UAS V2.1.4.500306
Content-Type: Application/MANSCDP+xml
Content-Length:   180

<?xml version="1.0" encoding="GB2312" ?>
<Notify>
    <CmdType>Keepalive</CmdType>
    <SN>204</SN>
    <DeviceID>34020000001320000001</DeviceID>
    <Status>OK</Status>
</Notify>`)

	sip := NewSIP()
	at.Nil(sip.ParseBytes(raw))
	at.Equal(SIPMethodMessage, sip.Method)
	at.Equal("34020000001320000001", sip.From.User)
	at.Equal("34020000002000000001", sip.To.User)
	at.Equal("192.168.1.102", sip.Via.Host)
	at.Equal("MESSAGE", sip.Cseq.Method)
	at.False(sip.IsResponse)

	at.Len(sip.Payload(), 180)
}

func TestSIP_Bytes(t *testing.T) {
	// at := assert.New(t)
	//
	// sip := NewSIP()
	// sip.Method
}
