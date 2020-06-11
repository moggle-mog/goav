// sips From https://github.com/google/gopacket/blob/master/layers/sip.go
package sips

import (
	"fmt"
	"strings"
)

// SIPVersion defines the different versions of the SIP Protocol
type SIPVersion uint8

// Represents all the versions of SIP protocol
const (
	SIPVersion1 SIPVersion = 1
	SIPVersion2 SIPVersion = 2
)

func (sv SIPVersion) String() string {
	switch sv {
	default:
		// Defaulting to SIP/2.0
		return "SIP/2.0"
	case SIPVersion1:
		return "SIP/1.0"
	case SIPVersion2:
		return "SIP/2.0"
	}
}

// GetSIPVersion is used to get SIP version constant
func GetSIPVersion(version string) (SIPVersion, error) {
	switch strings.ToUpper(version) {
	case "SIP/1.0":
		return SIPVersion1, nil
	case "SIP/2.0":
		return SIPVersion2, nil
	default:
		return 0, fmt.Errorf("unknown SIP version: '%s'", version)
	}
}
