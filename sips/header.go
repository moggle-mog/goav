// sips ref https://github.com/google/gopacket/blob/master/layers/sip.go
package sips

import (
	"bytes"
	"strconv"
	"strings"
)

// Here is a correspondence between long header names and short
// as defined in rfc3261 in section 20
var compactSipHeadersCorrespondence = map[string]string{
	"accept-contact":      "a",
	"allow-events":        "u",
	"call-id":             "i",
	"contact":             "m",
	"content-encoding":    "e",
	"content-length":      "l",
	"content-type":        "c",
	"event":               "o",
	"from":                "f",
	"identity":            "y",
	"refer-to":            "r",
	"referred-by":         "b",
	"reject-contact":      "j",
	"request-disposition": "d",
	"session-expires":     "x",
	"subject":             "s",
	"supported":           "k",
	"to":                  "t",
	"via":                 "v",
}

// SIPHeader sip's header
type SIPHeader map[string][]string

// NewSIPHeader returns a new sip header struct
func NewSIPHeader() SIPHeader {
	return make(SIPHeader)
}

// GetAllHeaders will return the full headers of the
// current SIP packets in a map[string][]string
func (s SIPHeader) GetAllHeaders() map[string][]string {
	return s
}

// GetHeader will return all the headers with
// the specified name.
func (s SIPHeader) GetHeader(headerName string) []string {
	headerName = strings.ToLower(headerName)
	h := make([]string, 0)
	if _, ok := s[headerName]; ok {
		if len(s[headerName]) > 0 {
			return s[headerName]
		} else if len(s[compactSipHeadersCorrespondence[headerName]]) > 0 {
			return s[compactSipHeadersCorrespondence[headerName]]
		}
	}
	return h
}

// GetFirstHeader will return the first header with
// the specified name. If the current SIP packet has multiple
// headers with the same name, it returns the first.
func (s SIPHeader) GetFirstHeader(headerName string) string {
	headerName = strings.ToLower(headerName)
	if _, ok := s[headerName]; ok {
		if len(s[headerName]) > 0 {
			return s[headerName][0]
		} else if len(s[compactSipHeadersCorrespondence[headerName]]) > 0 {
			return s[compactSipHeadersCorrespondence[headerName]][0]
		}
	}
	return ""
}

// SetFirstHeader sets header with only one value
func (s SIPHeader) SetFirstHeader(headerName, header string) {
	headerName = strings.ToLower(headerName)
	s[headerName] = []string{header}
}

// Bytes saves sip header into bytes
// The table below lists the header fields currently defined for the
// Session Initiation Protocol (SIP) [RFC3261]. Some headers have
// single-letter compact forms (Section 7.3 of RFC 3261). Header field
// names are case-insensitive.
//
// ref: https://www.iana.org/assignments/sip-parameters/sip-parameters.xhtml
//
func (s SIPHeader) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	for k, v := range s {

		buf.WriteString(k)
		buf.WriteString(":")
		for _, v1 := range v {
			buf.WriteString(" ")
			buf.WriteString(v1)
			buf.WriteString("\r\n")
			continue
		}
	}

	return buf.Bytes()
}

func (s SIPHeader) String() string {
	return string(s.Bytes())
}

// GetSIPFrom gets from struct
func (s SIPHeader) GetSIPFrom() *SIPUser {
	return NewSIPUser(s.GetFirstHeader("from"))
}

// GetSIPTo gets to struct
func (s SIPHeader) GetSIPTo() *SIPUser {
	return NewSIPUser(s.GetFirstHeader("to"))
}

// GetSIPContact gets contact struct
func (s SIPHeader) GetSIPContact() *SIPUser {
	return NewSIPUser(s.GetFirstHeader("contact"))
}

// GetSIPCseq gets cseq struct
func (s SIPHeader) GetSIPCseq() *SIPCseq {
	return NewSIPCseq(s.GetFirstHeader("cseq"))
}

// GetSIPVia gets via struct
func (s SIPHeader) GetSIPVia() *SIPVia {
	return NewSIPVia(s.GetFirstHeader("via"))
}

// GetContentLength gets content-length and cast to int
func (s SIPHeader) GetContentLength() (int, error) {
	return strconv.Atoi(s.GetFirstHeader("content-length"))
}

// GetCallID gets Call-ID string
func (s SIPHeader) GetCallID() string {
	return s.GetFirstHeader("call-id")
}

// GetAuthorization gets Authorization string
func (s SIPHeader) GetAuthorization() *SIPAuth {
	return NewSIPAuth(s.GetFirstHeader("authorization"))
}

// GetExpires gets Expires string
func (s SIPHeader) GetExpires() string {
	return s.GetFirstHeader("expires")
}

// CopyFrom copies from src header with header name
func (s SIPHeader) CopyFrom(src SIPHeader, names ...string) {
	for _, v := range names {
		s.SetFirstHeader(v, src.GetFirstHeader(v))
	}
}

// SetContentLength sets Content-Length header
func (s SIPHeader) SetContentLength(size int) {
	s.SetFirstHeader("content-length", strconv.Itoa(size))
}

// SetCallID sets Call-ID header
func (s SIPHeader) SetCallID(callID string) {
	s.SetFirstHeader("call-id", callID)
}

// SetMaxForwards sets Max-Forwards header
func (s SIPHeader) SetMaxForwards(max int) {
	s.SetFirstHeader("max-forwards", strconv.Itoa(max))
}

// SetReasonPhrase sets Reason-Phrase header
func (s SIPHeader) SetReasonPhrase(reason string) {
	s.SetFirstHeader("reason-phrase", reason)
}

// SetSupported sets Supported header
func (s SIPHeader) SetSupported(supported string) {
	s.SetFirstHeader("supported", supported)
}

// SetSubject sets Subject header
func (s SIPHeader) SetSubject(subject string) {
	s.SetFirstHeader("subject", subject)
}

// SetContentType sets Content-Type header
func (s SIPHeader) SetContentType(ct string) {
	s.SetFirstHeader("content-type", ct)
}

// SetAllow sets Allow header
func (s SIPHeader) SetAllow(allow string) {
	s.SetFirstHeader("allow", allow)
}

// SetSIPFrom sets From header
func (s SIPHeader) SetSIPFrom(from *SIPUser) {
	s.SetFirstHeader("from", from.String())
}

// SetSIPTo sets To header
func (s SIPHeader) SetSIPTo(to *SIPUser) {
	s.SetFirstHeader("to", to.String())
}

// SetSIPContact sets Contact header
func (s SIPHeader) SetSIPContact(contact *SIPUser) {
	s.SetFirstHeader("contact", contact.String())
}

// SetSIPCseq sets CSeq header
func (s SIPHeader) SetSIPCseq(cseq *SIPCseq) {
	s.SetFirstHeader("cseq", cseq.String())
}

// SetSIPVia sets Via header
func (s SIPHeader) SetSIPVia(via *SIPVia) {
	s.SetFirstHeader("via", via.String())
}
