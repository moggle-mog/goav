// sips ref
// https://github.com/marv2097/siprocket/blob/master/sipFrom.go
// https://github.com/marv2097/siprocket/blob/master/sipTo.go
// https://github.com/marv2097/siprocket/blob/master/sipContact.go
package sips

import (
	"bytes"
	"strconv"
	"strings"
)

// SIPUser a single line that is in the format of a from or to line
//
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 8.1.1.2 To
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 8.1.1.3 From
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 8.1.1.8 Contact
//
type SIPUser struct {
	URIType string // Type of URI sip, sips, tel etc
	Name    string // Name portion of URI
	User    string // User part
	Host    string // Host part
	Port    string // Port number
	Args    Args   // Args of line
	Src     string // Full source
}

// NewSIPUser parses URI into SIPUser struct
//
// Examples of user line of SIP Protocol :
//
// From/To: "Bob" <sips:bob@biloxi.com> ;tag=a48s
// From/To: sip:+12125551212@phone2net.com;tag=887s
// From/To: Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8
// From/To: <sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856

// Contact: "Mr. Watson" <sip:watson@worcester.bell-telephone.com>
//  ;q=0.7; expires=3600,
//  "Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1
// Contact: <sips:bob@192.0.2.4>;expires=60
// Contact: <sip:34020000001320000001@192.168.1.102:5060>
//
func NewSIPUser(uri string) *SIPUser {

	su := &SIPUser{
		URIType: "sip",
		Args:    make(Args),
		Src:     uri,
	}
	pos, state := 0, _FieldBase

	var name, user, host, port []byte

	for pos < len(uri) && uri[pos] != ';' {

		if uri[pos] == '>' || (uri[pos] == ' ' && state != _FieldName) {
			state = _FieldBase
			pos++
			continue
		}

		switch state {
		case _FieldBase:

			str := Substring(uri, pos, pos+4)
			if str == "sip:" || str == "tel:" {
				state = _FieldUser
				su.URIType = uri[pos : pos+3]
				pos = pos + 4
				continue
			}

			str = Substring(uri, pos, pos+5)
			if str == "sips:" {
				state = _FieldUser
				su.URIType = uri[pos : pos+4]
				pos = pos + 5
				continue
			}

			if uri[pos] == '"' || su.URIType == "" {
				state = _FieldName
				continue
			}
		case _FieldName:

			if uri[pos] == '"' || uri[pos] == '<' {
				state = _FieldBase
				pos++
				continue
			}
			name = append(name, uri[pos])
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

	su.Name = string(name)
	su.User = string(user)
	su.Host = string(host)
	su.Port = string(port)

	if len(su.Name) > 0 {
		su.Name = strings.TrimSpace(su.Name)
	}

	if pos < len(uri) {
		su.Args = ParseArgs(uri[pos:])
	}

	return su
}

// Bytes package SIPUser struct into slice
func (s *SIPUser) Bytes() []byte {

	buf := bytes.NewBuffer(nil)

	if len(s.Name) > 0 {

		if s.Name == "Anonymous" {
			buf.WriteString(s.Name)
		} else {
			buf.WriteString(strconv.Quote(s.Name))
		}

		buf.WriteString(" ")
	}

	buf.WriteString("<")
	buf.WriteString(s.URIType)
	buf.WriteString(":")
	buf.WriteString(s.User)
	buf.WriteString("@")
	buf.WriteString(s.Host)

	if len(s.Port) > 0 {

		buf.WriteString(":")
		buf.WriteString(s.Port)
	}

	buf.WriteString(">")

	if len(s.Args) > 0 {
		buf.WriteString(s.Args.SemicolonString())
	}

	s.Src = buf.String()

	return buf.Bytes()
}

func (s *SIPUser) String() string {
	return string(s.Bytes())
}

// ToRequest transform SIPUser struct to SIPRequest struct
func (s *SIPUser) ToRequest() *SIPRequest {

	sip := NewSIPRequest("")
	sip.URIType = s.URIType
	sip.User = s.User
	sip.Host = s.Host
	sip.Port = s.Port
	sip.Args = s.Args
	sip.Src = s.Src

	return sip
}
