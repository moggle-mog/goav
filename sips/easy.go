package sips

// MakeRequest make SIP struct to request
func MakeRequest(method SIPMethod, request *SIPRequest, header SIPHeader, body []byte) *SIP {

	sip := NewSIP()

	sip.Method = method
	sip.Request = request
	sip.Version = SIPVersion2

	if len(header) > 0 {
		sip.SetHeader(header)
	}
	if len(body) > 0 {
		sip.SetBody(body)
	}

	return sip
}

// MakeResponse make SIP struct to respond
func MakeResponse(code SIPStatus, header SIPHeader, body []byte) *SIP {

	sip := NewSIP()
	sip.IsResponse = true

	sip.Version = SIPVersion2
	sip.ResponseCode = code

	if len(header) > 0 {
		sip.SetHeader(header)
	}
	if len(body) > 0 {
		sip.SetBody(body)
	}

	return sip
}
