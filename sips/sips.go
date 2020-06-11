// sips ref https://github.com/google/gopacket/blob/master/layers/sip.go
package sips

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// SIP object will contains information about decoded SIP packet.
// -> The SIP Version
// -> The SIP Headers (in a map[string][]string because of multiple headers with the same name
// -> The SIP Method
// -> The SIP Response code (if it's a response)
// -> The SIP Status line (if it's a response)
// You can easily know the type of the packet with the IsResponse boolean
//
type SIP struct {
	BaseLayer

	// Base information
	Version SIPVersion
	Headers SIPHeader

	// Common Header
	From    *SIPUser
	To      *SIPUser
	Contact *SIPUser
	Cseq    *SIPCseq
	Via     *SIPVia

	// Request
	Method  SIPMethod
	Request *SIPRequest

	// Response
	IsResponse     bool
	ResponseCode   SIPStatus
	ResponseStatus string

	// Private fields
	lastHeaderParsed string
}

// NewSIP instantiates a new empty SIP object
func NewSIP() *SIP {
	s := new(SIP)
	s.Headers = make(map[string][]string)
	return s
}

// Payload returns the base layer payload
func (s *SIP) Payload() []byte {
	return s.BaseLayer.LayerPayload()
}

func (s *SIP) fillHeader() {

	s.Contact = s.Headers.GetSIPContact()
	s.Cseq = s.Headers.GetSIPCseq()
	s.From = s.Headers.GetSIPFrom()
	s.To = s.Headers.GetSIPTo()
	s.Via = s.Headers.GetSIPVia()
}

func (s *SIP) fillBody(r *bufio.Reader, supplement int) ([]byte, error) {

	var body []byte

	size, err := s.Headers.GetContentLength()
	if err != nil {
		return nil, err
	}
	size -= supplement

	for size > 0 {

		// read one block
		buf := make([]byte, 4*1024)

		n, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			if n == 0 {
				break
			}

			// save remain byte
			size -= n
			body = append(body, buf[:n]...)
			break
		}

		size -= n
		body = append(body, buf[:n]...)
	}

	return body, nil
}

// Parse parses the slices from reader into the SIP struct.
func (s *SIP) Parse(r *bufio.Reader) error {

	// Init some vars for parsing follow-up
	var countLines int
	var offset int
	var buffer []byte

	for {

		// Read next line
		line, err := r.ReadSlice(byte('\n'))
		if err != nil {
			if err != io.EOF {
				return err
			}

			// ReadSlice返回EOF时line可能不为空，因此仍需收尾
			if len(line) == 0 {
				break
			}
		}

		buffer = append(buffer, line...)
		offset += len(line)

		// Trim the new line delimiters
		trim := bytes.Trim(line, "\r\n")

		// Empty line, we hit Body
		if len(trim) == 0 {
			offset -= len(line)
			break
		}
		line = trim

		// First line is the SIP request/response line
		// Other lines are headers
		if countLines == 0 {
			err = s.ParseFirstLine(line)
			if err != nil {
				return err
			}

		} else {
			s.ParseHeader(line)
		}

		countLines++
	}

	s.fillHeader()

	// The size of the message-body does not include the CRLF separating
	// header fields and body.  Any Content-Length greater than or equal to
	// zero is a valid value.  If no body is present in a message, then the
	// Content-Length header field value MUST be set to zero.
	body, err := s.fillBody(r, len(buffer)-offset)
	if err != nil {
		return err
	}
	buffer = append(buffer, body...)

	s.BaseLayer = BaseLayer{Contents: buffer[:offset], Payload: buffer[offset:]}
	return nil
}

// ParseBytes parses the slice into the SIP struct.
func (s *SIP) ParseBytes(buf []byte) error {
	return s.Parse(bufio.NewReader(bytes.NewReader(buf)))
}

// ParseHeader will parse a SIP Header
// SIP Headers are quite simple, there are colon separated name and value
// Headers can be spread over multiple lines
//
// Examples of header :
//
//  CSeq: 1 REGISTER
//  Via: SIP/2.0/UDP there.com:5060
//  Authorization:Digest username="UserB",
//	  realm="MCI WorldCom SIP",
//    nonce="1cec4341ae6cbe5a359ea9c8e88df84f", opaque="",
//    uri="sip:ss2.wcom.com", response="71ba27c64bd01de719686aa4590d5824"
//
func (s *SIP) ParseHeader(header []byte) {

	// Ignore empty headers
	if len(header) == 0 {
		return
	}

	// Check if this is the following of last header
	// RFC 3261 - 7.3.1 - Header Field Format specify that following lines of
	// multiline headers must begin by SP or TAB
	if header[0] == '\t' || header[0] == ' ' {

		header = bytes.TrimSpace(header)
		s.Headers[s.lastHeaderParsed][len(s.Headers[s.lastHeaderParsed])-1] += fmt.Sprintf(" %s", string(header))
		return
	}

	// Find the ':' to separate header name and value
	index := bytes.Index(header, []byte(":"))
	if index >= 0 {

		headerName := strings.ToLower(string(bytes.Trim(header[:index], " ")))
		headerValue := string(bytes.Trim(header[index+1:], " "))

		// Add header to object
		s.Headers[headerName] = append(s.Headers[headerName], headerValue)
		s.lastHeaderParsed = headerName
	}
}

// Encode package the SIP struct into slice
func (s *SIP) Bytes() []byte {

	buf := bytes.NewBuffer(nil)

	// package request
	if s.IsResponse {

		// --> Response
		buf.WriteString(s.Version.String())
		buf.WriteString(" ")
		buf.WriteString(s.ResponseCode.String())
		buf.WriteString(" ")
		buf.WriteString(s.ResponseCode.Text())
	} else {

		// --> Request
		buf.WriteString(s.Method.String())
		buf.WriteString(" ")
		buf.WriteString(s.Request.String())
		buf.WriteString(" ")
		buf.WriteString(s.Version.String())
	}

	buf.WriteString("\r\n")

	// add headers - Content-Length
	s.Headers.SetContentLength(len(s.Payload()))

	// package header
	buf.Write(s.Headers.Bytes())

	// package body
	if len(s.Payload()) > 0 {
		buf.WriteString("\r\n")
		buf.Write(s.Payload())
	}

	return buf.Bytes()
}

func (s *SIP) String() string {
	return string(s.Bytes())
}

// SetHeader sets header into SIPHeader struct
func (s *SIP) SetHeader(header SIPHeader) {
	s.Headers = header
}

// SetBody sets body into BaseLayer struct
func (s *SIP) SetBody(body []byte) {
	s.BaseLayer.Payload = body
}
