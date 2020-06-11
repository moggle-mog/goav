// sips ref
// https://github.com/google/gopacket/blob/master/layers/sip.go
// https://github.com/marv2097/siprocket/blob/master/sipRequestLine.go
package sips

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// SIPRequest the Request-URI is a SIP or SIPS URI as described in
// Section 19.1 or a general URI (RFC 2396 [5]).  It indicates
// the user or service to which this request is being addressed.
// The Request-URI MUST NOT contain unescaped spaces or control
// characters and MUST NOT be enclosed in "<>".
//
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt 7.1
//
// Examples of first line of SIP Protocol :
//
// Request：     INVITE sip:01798300765@87.252.61.202:5060;user=phone SIP/2.0
// Response：    SIP/2.0 200 OK
//
type SIPRequest struct {
	URIType string // Type of URI sip, sips, tel etc.
	User    string // User 01798300765 e.g.
	Host    string // Host 87.252.61.202 e.g.
	Port    string // Port 5060 e.g.
	Args    Args   // Args of line
	Src     string // Source
}

// NewSIPRequest parses uri into request struct
//
// URI Example:
//
// sip:01798300765@87.252.61.202:5060;user=phone
//
func NewSIPRequest(uri string) *SIPRequest {

	sr := &SIPRequest{
		URIType: "sip",
		Args:    make(Args),
		Src:     uri,
	}
	pos, state := 0, _FieldBase

	var user, host, port []byte

	for pos < len(uri) && uri[pos] != ';' {

		switch state {
		case _FieldBase:

			str := Substring(uri, pos, pos+4)
			if str == "sip:" || str == "tel:" {
				state = _FieldUser
				sr.URIType = uri[pos : pos+3]
				pos = pos + 4
				continue
			}

			str = Substring(uri, pos, pos+5)
			if str == "sips:" {
				state = _FieldUser
				sr.URIType = uri[pos : pos+4]
				pos = pos + 5
				continue
			}
		case _FieldUser:

			if uri[pos] == '@' {
				state = _FieldHost
				pos++
				continue
			}
			user = append(user, uri[pos])
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

	sr.User = string(user)
	sr.Host = string(host)
	sr.Port = string(port)

	if pos < len(uri) {
		sr.Args = ParseArgs(uri[pos:])
	}

	return sr
}

// Bytes package Request struct into slice
func (s *SIPRequest) Bytes() []byte {

	buf := bytes.NewBuffer(nil)

	// sip:
	buf.WriteString(s.URIType)
	buf.WriteString(":")

	// user@host
	buf.WriteString(s.User)
	buf.WriteString("@")
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

func (s *SIPRequest) String() string {
	return string(s.Bytes())
}

// ToUser transform SIPRequest struct to SIPUser struct
func (s *SIPRequest) ToUser() *SIPUser {

	sip := NewSIPUser("")
	sip.URIType = s.URIType
	sip.User = s.User
	sip.Host = s.Host
	sip.Port = s.Port
	sip.Args = s.Args
	sip.Src = s.Src

	return sip
}

// ParseFirstLine will compute the first line of a SIP packet.
// The first line will tell us if it's a request or a response.
//
// Examples of first line of SIP Protocol :
//
// 	Request 	: INVITE bob@example.com SIP/2.0
// 	Response 	: SIP/2.0 200 OK
// 	Response	: SIP/2.0 501 Not Implemented
//
func (s *SIP) ParseFirstLine(firstLine []byte) error {

	var err error

	// Splits line by space
	splits := strings.SplitN(string(firstLine), " ", 3)

	// We must have at least 3 parts
	if len(splits) < 3 {
		return fmt.Errorf("invalid first SIP line: '%s'", string(firstLine))
	}

	// Determine the SIP packet type
	if strings.HasPrefix(splits[0], "SIP") {

		// --> Response
		s.IsResponse = true

		// Validate SIP Version
		s.Version, err = GetSIPVersion(splits[0])
		if err != nil {
			return err
		}

		// Compute code
		code, err := strconv.Atoi(splits[1])
		if err != nil {
			return err
		}
		s.ResponseCode = SIPStatus(code)

		// Compute status line
		s.ResponseStatus = splits[2]

	} else {

		// --> Request

		// Validate method
		s.Method, err = GetSIPMethod(splits[0])
		if err != nil {
			return err
		}

		// request URI
		s.Request = NewSIPRequest(splits[1])

		// Validate SIP Version
		s.Version, err = GetSIPVersion(splits[2])
		if err != nil {
			return err
		}

	}

	return nil
}
