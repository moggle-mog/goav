package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSIPHeader_Bytes(t *testing.T) {
	at := assert.New(t)

	sh := make(SIPHeader)
	sh.SetFirstHeader("Via", "SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK1135200957")
	at.Equal("via: SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK1135200957\r\n", sh.String())
	at.Equal("SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK1135200957", sh.GetFirstHeader("via"))

	sh = make(SIPHeader)
	sh.SetFirstHeader("from", "<sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856")
	at.Equal("from: <sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856\r\n", sh.String())

	sh = make(SIPHeader)
	sh.SetFirstHeader("To", "<sip:34020000002000000001@10.104.157.255:5060>")
	at.Equal("to: <sip:34020000002000000001@10.104.157.255:5060>\r\n", sh.String())

	sh = make(SIPHeader)
	sh.SetFirstHeader("Call-ID", "1808031968@192.168.1.102")
	at.Equal("call-id: 1808031968@192.168.1.102\r\n", sh.String())

	sh = make(SIPHeader)
	sh.SetFirstHeader("CSeq", "20 MESSAGE")
	at.Equal("cseq: 20 MESSAGE\r\n", sh.String())

	sh = make(SIPHeader)
	sh.SetFirstHeader("Max-Forwards", "70")
	at.Equal("max-forwards: 70\r\n", sh.String())

	sh = make(SIPHeader)
	sh["User-Agent"] = []string{
		"SIP UAS V2.1.4.500306",
	}
	at.Equal("User-Agent: SIP UAS V2.1.4.500306\r\n", sh.String())

	sh = make(SIPHeader)
	sh["Contact"] = []string{
		`"Mr. Watson" <sip:watson@worcester.bell-telephone.com>`,
		` ;q=0.7; expires=3600,`,
		` "Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1`,
	}
	at.Equal(`Contact: "Mr. Watson" <sip:watson@worcester.bell-telephone.com>`+"\r\n"+
		`  ;q=0.7; expires=3600,`+"\r\n"+
		`  "Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1`+"\r\n", sh.String())

	sh = make(SIPHeader)
	sh["Content-Type"] = []string{
		"Application/MANSCDP+xml",
	}
	at.Equal("Content-Type: Application/MANSCDP+xml\r\n", sh.String())

	sh = make(SIPHeader)
	sh["Content-Length"] = []string{
		"0",
	}
	at.Equal("Content-Length: 0\r\n", sh.String())
}
