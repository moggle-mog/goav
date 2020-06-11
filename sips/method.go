// sips From https://github.com/google/gopacket/blob/master/layers/sip.go
package sips

import (
	"fmt"
	"strings"
)

// SIPMethod defines the different methods of the SIP Protocol
// defined in the different RFC's
type SIPMethod uint16

// Here are all the SIP methods
const (
	SIPMethodNil       SIPMethod = 0
	SIPMethodInvite    SIPMethod = 1  // INVITE	[RFC3261]
	SIPMethodAck       SIPMethod = 2  // ACK	[RFC3261]
	SIPMethodBye       SIPMethod = 3  // BYE	[RFC3261]
	SIPMethodCancel    SIPMethod = 4  // CANCEL	[RFC3261]
	SIPMethodOptions   SIPMethod = 5  // OPTIONS	[RFC3261]
	SIPMethodRegister  SIPMethod = 6  // REGISTER	[RFC3261]
	SIPMethodPrack     SIPMethod = 7  // PRACK	[RFC3262]
	SIPMethodSubscribe SIPMethod = 8  // SUBSCRIBE	[RFC6665]
	SIPMethodNotify    SIPMethod = 9  // NOTIFY	[RFC6665]
	SIPMethodPublish   SIPMethod = 10 // PUBLISH	[RFC3903]
	SIPMethodInfo      SIPMethod = 11 // INFO	[RFC6086]
	SIPMethodRefer     SIPMethod = 12 // REFER	[RFC3515]
	SIPMethodMessage   SIPMethod = 13 // MESSAGE	[RFC3428]
	SIPMethodUpdate    SIPMethod = 14 // UPDATE	[RFC3311]
	SIPMethodPing      SIPMethod = 15 // PING	[https://tools.ietf.org/html/draft-fwmiller-ping-03]
)

func (sm SIPMethod) String() string {
	switch sm {
	default:
		return "Unknown method"
	case SIPMethodInvite:
		return "INVITE"
	case SIPMethodAck:
		return "ACK"
	case SIPMethodBye:
		return "BYE"
	case SIPMethodCancel:
		return "CANCEL"
	case SIPMethodOptions:
		return "OPTIONS"
	case SIPMethodRegister:
		return "REGISTER"
	case SIPMethodPrack:
		return "PRACK"
	case SIPMethodSubscribe:
		return "SUBSCRIBE"
	case SIPMethodNotify:
		return "NOTIFY"
	case SIPMethodPublish:
		return "PUBLISH"
	case SIPMethodInfo:
		return "INFO"
	case SIPMethodRefer:
		return "REFER"
	case SIPMethodMessage:
		return "MESSAGE"
	case SIPMethodUpdate:
		return "UPDATE"
	case SIPMethodPing:
		return "PING"
	}
}

// GetSIPMethod returns the constant of a SIP method
// from its string
func GetSIPMethod(method string) (SIPMethod, error) {
	switch strings.ToUpper(method) {
	case "INVITE":
		return SIPMethodInvite, nil
	case "ACK":
		return SIPMethodAck, nil
	case "BYE":
		return SIPMethodBye, nil
	case "CANCEL":
		return SIPMethodCancel, nil
	case "OPTIONS":
		return SIPMethodOptions, nil
	case "REGISTER":
		return SIPMethodRegister, nil
	case "PRACK":
		return SIPMethodPrack, nil
	case "SUBSCRIBE":
		return SIPMethodSubscribe, nil
	case "NOTIFY":
		return SIPMethodNotify, nil
	case "PUBLISH":
		return SIPMethodPublish, nil
	case "INFO":
		return SIPMethodInfo, nil
	case "REFER":
		return SIPMethodRefer, nil
	case "MESSAGE":
		return SIPMethodMessage, nil
	case "UPDATE":
		return SIPMethodUpdate, nil
	case "PING":
		return SIPMethodPing, nil
	default:
		return 0, fmt.Errorf("Unknown SIP method: '%s'", method)
	}
}
