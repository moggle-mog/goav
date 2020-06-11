package sips

import "strconv"

// SIPStatus status code
type SIPStatus int

// SIP response status codes.
const (
	// 100+
	StatusTrying               SIPStatus = 100
	StatusRinging              SIPStatus = 180
	StatusCallIsBeingForwarded SIPStatus = 181
	StatusQueued               SIPStatus = 182
	StatusSessionProgress      SIPStatus = 183

	// 200+
	StatusOK SIPStatus = 200

	// 300+
	StatusMultipleChoices    SIPStatus = 300
	StatusMovedPermanently   SIPStatus = 301
	StatusMovedTemporarily   SIPStatus = 302
	StatusUseProxy           SIPStatus = 305
	StatusAlternativeService SIPStatus = 380

	// 400+
	StatusBadRequest                  SIPStatus = 400
	StatusUnauthorized                SIPStatus = 401
	StatusPaymentRequired             SIPStatus = 402
	StatusForbidden                   SIPStatus = 403
	StatusNotFound                    SIPStatus = 404
	StatusMethodNotAllowed            SIPStatus = 405
	StatusNotAcceptable               SIPStatus = 406
	StatusProxyAuthenticationRequired SIPStatus = 407
	StatusRequestTimeout              SIPStatus = 408
	StatusGone                        SIPStatus = 410
	StatusRequestEntityTooLarge       SIPStatus = 413
	StatusRequestURITooLong           SIPStatus = 414
	StatusUnsupportedMediaType        SIPStatus = 415
	StatusUnsupportedURIScheme        SIPStatus = 416
	StatusBadExtension                SIPStatus = 420
	StatusExtensionRequired           SIPStatus = 421
	StatusIntervalTooBrief            SIPStatus = 423
	StatusNoResponse                  SIPStatus = 480
	StatusCallTransactionDoesNotExist SIPStatus = 481
	StatusLoopDetected                SIPStatus = 482
	StatusTooManyHops                 SIPStatus = 483
	StatusAddressIncomplete           SIPStatus = 484
	StatusAmbigious                   SIPStatus = 485
	StatusBusyHere                    SIPStatus = 486
	StatusRequestTerminated           SIPStatus = 487
	StatusNotAcceptableHere           SIPStatus = 488
	StatusRequestPending              SIPStatus = 491
	StatusUndecipherable              SIPStatus = 493

	// 500+
	StatusServerInternalError SIPStatus = 500
	StatusNotImplemented      SIPStatus = 501
	StatusBadGateway          SIPStatus = 502
	StatusServiceUnavailable  SIPStatus = 503
	StatusServerTimeout       SIPStatus = 504
	StatusVersionNotSupported SIPStatus = 505
	StatusMessageTooLarge     SIPStatus = 513

	// 600+
	StatusBusyEverywhere       SIPStatus = 600
	StatusDecline              SIPStatus = 603
	StatusDoesNotExistAnywhere SIPStatus = 604
	StatusUnacceptable         SIPStatus = 606
)

var statusTexts = map[SIPStatus]string{
	StatusTrying:                      "Trying",
	StatusRinging:                     "Ringing",
	StatusCallIsBeingForwarded:        "Call Is Being Forwarded",
	StatusQueued:                      "Queued",
	StatusSessionProgress:             "Session Progress",
	StatusOK:                          "OK",
	StatusMultipleChoices:             "Multiple Choices",
	StatusMovedPermanently:            "Moved Permanently",
	StatusMovedTemporarily:            "Moved Temporarily",
	StatusUseProxy:                    "Use Proxy",
	StatusAlternativeService:          "Alternative Service",
	StatusBadRequest:                  "Bad Request",
	StatusUnauthorized:                "Unauthorized",
	StatusPaymentRequired:             "Payment Required",
	StatusForbidden:                   "Forbidden",
	StatusNotFound:                    "Not Found",
	StatusMethodNotAllowed:            "Method Not Allowed",
	StatusNotAcceptable:               "Not Acceptable",
	StatusProxyAuthenticationRequired: "Proxy Authentication Required",
	StatusRequestTimeout:              "Request Timeout",
	StatusGone:                        "Gone",
	StatusRequestEntityTooLarge:       "Request Entity Too Large",
	StatusRequestURITooLong:           "Request-URI Too Long",
	StatusUnsupportedMediaType:        "Unsupported Media Type",
	StatusUnsupportedURIScheme:        "Unsupported URI Scheme",
	StatusBadExtension:                "Bad Extension",
	StatusExtensionRequired:           "Extension Required",
	StatusIntervalTooBrief:            "Interval Too Brief",
	StatusNoResponse:                  "No Response",
	StatusCallTransactionDoesNotExist: "Call/Transaction Does Not Exist",
	StatusLoopDetected:                "Loop Detected",
	StatusTooManyHops:                 "Too Many Hops",
	StatusAddressIncomplete:           "Address Incomplete",
	StatusAmbigious:                   "Ambiguous",
	StatusBusyHere:                    "Busy Here",
	StatusRequestTerminated:           "Request Terminated",
	StatusNotAcceptableHere:           "Not Acceptable Here",
	StatusRequestPending:              "Request Pending",
	StatusUndecipherable:              "Undecipherable",
	StatusServerInternalError:         "Server Internal Error",
	StatusNotImplemented:              "Not Implemented",
	StatusBadGateway:                  "Bad Gateway",
	StatusServiceUnavailable:          "Service Unavailable",
	StatusServerTimeout:               "Server Timeout",
	StatusVersionNotSupported:         "Version Not Supported",
	StatusMessageTooLarge:             "Message Too Large",
	StatusBusyEverywhere:              "Busy Everywhere",
	StatusDecline:                     "Decline",
	StatusDoesNotExistAnywhere:        "Does Not Exist Anywhere",
	StatusUnacceptable:                "Not Acceptable",
}

// Text return the human readable text representation of a status code.
func (s SIPStatus) Text() string {
	return statusTexts[s]
}

// Int return code
func (s SIPStatus) Int() int {
	return int(s)
}

// String return code string
func (s SIPStatus) String() string {
	return strconv.Itoa(s.Int())
}
