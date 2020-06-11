// sips parse sip header
// ref https://github.com/1lann/go-sip/blob/master/sipnet/header_args.go
package sips

import (
	"bytes"
	"strconv"
	"strings"
)

// Args represents the arguments which can be found in headers,
// and in other simple key value fields whose format is of a
// key=value with a delimiter.
type Args map[string]string

// NewArgs returns new Args
func NewArgs() Args {
	return make(Args)
}

// ParseArgsPairs extracts key/value pairs from comma, semicolon, or new line
// separated values.
//
// Lifted from https://code.google.com/p/gorilla/source/browse/http/parser/parser.go
func ParseArgsPairs(value string) Args {
	m := make(Args)
	for _, pair := range m.ParseList(strings.TrimSpace(value)) {
		if i := strings.Index(pair, "="); i < 0 {
			m[pair] = ""
		} else {
			v := pair[i+1:]
			if v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			m[pair[:i]] = v
		}
	}
	return m
}

// ParseArgs parses header arguments from a full header.
func ParseArgs(str string) Args {
	argLocation := strings.Index(str, ";")
	if argLocation < 0 {
		return make(Args)
	}

	return ParseArgsPairs(str[argLocation+1:])
}

// ParseList parses a comma, semicolon, or new line separated list of values
// and returns list elements.
//
// Lifted from https://code.google.com/p/gorilla/source/browse/http/parser/parser.go
// which was ported from urllib2.parse_http_list, from the Python
// standard library.
func (a Args) ParseList(value string) []string {
	var list []string
	var escape, quote bool
	b := new(bytes.Buffer)
	for _, r := range value {
		switch {
		case escape:
			b.WriteRune(r)
			escape = false
		case quote:
			if r == '\\' {
				escape = true
			} else {
				if r == '"' {
					quote = false
				}
				b.WriteRune(r)
			}
		case r == ',' || r == ';' || r == '\n':
			list = append(list, strings.TrimSpace(b.String()))
			b.Reset()
		case r == '"':
			quote = true
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	// Append last part.
	if s := b.String(); s != "" {
		list = append(list, strings.TrimSpace(s))
	}
	return list
}

// Del deletes the key and its value from the header arguments. Deleting a non-existent
// key is a no-op.
func (a Args) Del(key string) {
	delete(a, key)
}

// Get returns the value at a given key. It returns an empty string if the
// key does not exist.
func (a Args) Get(key string) string {
	return a[key]
}

// Set sets a header argument key with a value.
func (a Args) Set(key, value string) {
	a[key] = value
}

// SemicolonString returns the header arguments as a semicolon
// separated unquoted strings with a leading semicolon.
func (a Args) SemicolonString() string {
	var result string
	for key, value := range a {
		if value == "" {
			result += ";" + key
		} else {
			result += ";" + key + "=" + value
		}
	}
	return result
}

// CommaString returns the header arguments as a comma and space
// separated string.
func (a Args) CommaString() string {
	if len(a) == 0 {
		return ""
	}

	var result string
	for key, value := range a {
		result += key + "=" + strconv.Quote(value) + ", "
	}
	return result[:len(result)-2]
}

// CRLFString returns the header arguments as a CRLF separated string.
func (a Args) CRLFString() string {
	if len(a) == 0 {
		return ""
	}

	var result string
	for key, value := range a {
		result += key + "=" + value + "\r\n"
	}
	return result
}
