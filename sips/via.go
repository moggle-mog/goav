// sips ref https://github.com/marv2097/siprocket/blob/master/sipVia.go
package sips

import "bytes"

// SIPVia a single line that is in the format of a Via line
//
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 8.1.1.7 Via
//
// The Via header field indicates the transport used for the transaction
// and identifies the location where the response is to be sent.  A Via
// header field value is added only after the transport that will be
// used to reach the next hop has been selected (which may involve the
// usage of the procedures in [4]).
//
type SIPVia struct {
	Trans string // Type of Transport udp, tcp, tls, sctp etc
	Host  string // Host part
	Port  string // Port number
	Args  Args   // Args of line
	Src   string
}

// NewSIPVia parses URI into SIPVia struct
//
// Examples of user line of SIP Protocol :
//
// Via: SIP/2.0/UDP 192.168.1.102:5060;rport;branch=z9hG4bK352699498
//
func NewSIPVia(uri string) *SIPVia {

	sv := &SIPVia{
		Args: make(Args),
		Src:  uri,
	}
	pos, state := 0, _FieldBase

	var host, port []byte

	for pos < len(uri) && uri[pos] != ';' {

		if uri[pos] == ' ' {
			pos++
			continue
		}

		switch state {
		case _FieldBase:

			// When the UAC creates a request, it MUST insert a Via into that
			// request.  The protocol name and protocol version in the header field
			// MUST be SIP and 2.0, respectively.  The Via header field value MUST
			// contain a branch parameter.  This parameter is used to identify the
			// transaction created by that request.  This parameter is used by both
			// the client and the server.
			str := Substring(uri, pos, pos+8)
			if str == "SIP/2.0/" {

				state = _FieldHost
				pos = pos + 8

				str = Substring(uri, pos, pos+3)
				if str == "UDP" || str == "TCP" || str == "TLS" {
					sv.Trans = str
					pos = pos + 3
					continue
				}

				str = Substring(uri, pos, pos+4)
				if str == "SCTP" {
					sv.Trans = str
					pos = pos + 4
					continue
				}
			}
		case _FieldHost:

			if uri[pos] == ':' {
				state = _FieldPort
				pos++
				continue
			}
			host = append(host, uri[pos])
		case _FieldPort:

			port = append(port, uri[pos])
		}

		pos++
	}

	sv.Host = string(host)
	sv.Port = string(port)

	if pos < len(uri) {
		sv.Args = ParseArgs(uri[pos:])
	}

	return sv
}

// Bytes package SIPVia struct into slice
func (s *SIPVia) Bytes() []byte {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(SIPVersion2.String())
	buf.WriteString("/")
	buf.WriteString(s.Trans)
	buf.WriteString(" ")
	buf.WriteString(s.Host)

	if len(s.Port) > 0 {
		buf.WriteString(":")
		buf.WriteString(s.Port)
	}
	if len(s.Args) > 0 {
		buf.WriteString(s.Args.SemicolonString())
	}

	s.Src = buf.String()

	return buf.Bytes()
}

func (s *SIPVia) String() string {
	return string(s.Bytes())
}
