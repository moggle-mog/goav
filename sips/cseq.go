package sips

import "bytes"

// SIPCseq The CSeq header field serves as a way to identify and order
// transactions.  It consists of a sequence number and a method.  The
// method MUST match that of the request.
//
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 8.1.1.5 CSeq
//
type SIPCseq struct {
	ID     string // Cseq ID
	Method string // Cseq Method
	Src    string
}

// NewSIPCseq parses URI into SIPCseq struct
//
// Examples of cseq line of SIP Protocol :
//
// CSeq: 4711 INVITE
//
func NewSIPCseq(uri string) *SIPCseq {

	sc := &SIPCseq{
		Src: uri,
	}
	pos, state := 0, _FieldID

	var id, method []byte

	for pos < len(uri) {

		switch state {
		case _FieldID:

			if uri[pos] == ' ' {
				state = _FieldMethod
				pos++
				continue
			}
			id = append(id, uri[pos])
		case _FieldMethod:

			method = append(method, uri[pos])
		}

		pos++
	}

	sc.ID = string(id)
	sc.Method = string(method)

	return sc
}

// Bytes package SIPUser struct into slice
func (s *SIPCseq) Bytes() []byte {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(s.ID)
	buf.WriteString(" ")
	buf.WriteString(s.Method)

	s.Src = buf.String()

	return buf.Bytes()
}

func (s *SIPCseq) String() string {
	return string(s.Bytes())
}
