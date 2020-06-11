package sips

import "strings"

// SIPAuth a single line that is in the format of a from or to line
//
// RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt - 20.7 Authorization
type SIPAuth struct {
	Args Args
	Src  string
}

// NewSIPAuth parses line into SIPAuth struct
//
// Examples of user line of SIP Protocol :
//
// Authorization: Digest username="Alice", realm="atlanta.com",
//  nonce="84a4cc6f3082121f32b42a2187831a9e",
//  response="7587245234b3434cc3412213e5f113a5432"
//
func NewSIPAuth(src string) *SIPAuth {

	sa := &SIPAuth{
		Args: NewArgs(),
		Src:  src,
	}

	if len(src) >= 8 && strings.ToLower(src[:7]) == "digest " {

		sa.Args = ParseArgsPairs(src[7:])
	}

	return sa
}
