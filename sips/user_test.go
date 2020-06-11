package sips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSIPUser(t *testing.T) {
	at := assert.New(t)

	u := NewSIPUser(`"Bob" <sips:bob@biloxi.com> ;tag=a48s`)
	at.Equal("Bob", u.Name)
	at.Equal("bob", u.User)
	at.Equal("biloxi.com", u.Host)
	at.Equal("", u.Port)
	at.Equal("a48s", u.Args["tag"])
	//
	u = NewSIPUser(`"Bob bb" <sips:bob@biloxi.com> ;tag=a48s`)
	at.Equal("Bob bb", u.Name)
	at.Equal("bob", u.User)
	at.Equal("biloxi.com", u.Host)
	at.Equal("", u.Port)
	at.Equal("a48s", u.Args["tag"])

	u = NewSIPUser("sip:+12125551212@phone2net.com;tag=887s")
	at.Equal("", u.Name)
	at.Equal("+12125551212", u.User)
	at.Equal("phone2net.com", u.Host)
	at.Equal("", u.Port)
	at.Equal("887s", u.Args["tag"])

	u = NewSIPUser("Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8")
	at.Equal("Anonymous", u.Name)
	at.Equal("c8oqz84zk7z", u.User)
	at.Equal("privacy.org", u.Host)
	at.Equal("", u.Port)
	at.Equal("hyh8", u.Args["tag"])

	u = NewSIPUser("<sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856")
	at.Equal("", u.Name)
	at.Equal("34020000001320000001", u.User)
	at.Equal("192.168.1.102", u.Host)
	at.Equal("5060", u.Port)
	at.Equal("1086912856", u.Args["tag"])

	u = NewSIPUser(" <sip:34020000001320000001@192.168.1.102:5060> ;tag=1086912856 ")
	at.Equal("", u.Name)
	at.Equal("34020000001320000001", u.User)
	at.Equal("192.168.1.102", u.Host)
	at.Equal("5060", u.Port)
	at.Equal("1086912856", u.Args["tag"])
}

func TestNewSIPUser2(t *testing.T) {
	at := assert.New(t)

	sc := NewSIPUser("<sip:34020000001320000001@192.168.1.102:5060>")
	at.Equal("sip", sc.URIType)
	at.Equal("192.168.1.102", sc.Host)
	at.Equal("5060", sc.Port)
	at.Len(sc.Args, 0)

	sc = NewSIPUser("<sips:bob@192.0.2.4>;expires=60")
	at.Equal("sips", sc.URIType)
	at.Equal("", sc.Name)
	at.Equal("bob", sc.User)
	at.Equal("192.0.2.4", sc.Host)
	at.Equal("60", sc.Args["expires"])

	sc = NewSIPUser(`"Mr. Watson" <sip:watson@worcester.bell-telephone.com> ;q=0.7; expires=3600`)
	at.Equal("sip", sc.URIType)
	at.Equal("Mr. Watson", sc.Name)
	at.Equal("watson", sc.User)
	at.Equal("worcester.bell-telephone.com", sc.Host)
	at.Equal("0.7", sc.Args["q"])
	at.Equal("3600", sc.Args["expires"])

	sc = NewSIPUser(`"Mr. Watson" <sip:watson@worcester.bell-telephone.com> ;q=0.7; expires=3600, "Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1`)
	at.Equal("Mr. Watson", sc.Name)
	at.Equal("sip", sc.URIType)
	at.Equal("watson", sc.User)
	at.Equal("worcester.bell-telephone.com", sc.Host)
	at.Equal("", sc.Port)
	at.Equal("0.1", sc.Args["q"])
	at.Equal("", sc.Args[`"Mr. Watson" <mailto:watson@bell-telephone.com>`])
	at.Equal("3600", sc.Args["expires"])
}

func TestNewSIPUser1(t *testing.T) {
	at := assert.New(t)

	u := NewSIPUser(`sip:+12125551212@phone2net.com;tag=887s`)
	at.Equal("", u.Name)
	at.Equal("+12125551212", u.User)
	at.Equal("phone2net.com", u.Host)
	at.Equal("", u.Port)
	at.Equal("887s", u.Args["tag"])
}

func TestSIPUser_Bytes(t *testing.T) {
	at := assert.New(t)

	// "Bob" <sips:bob@biloxi.com> ;tag=a48s
	sip := NewSIPUser("")
	sip.Name = "Bob"
	sip.URIType = "sips"
	sip.User = "bob"
	sip.Host = "biloxi.com"
	sip.Args.Set("tag", "a48s")
	at.Equal(`"Bob" <sips:bob@biloxi.com>;tag=a48s`, sip.String())

	// sip:+12125551212@phone2net.com;tag=887s
	sip = NewSIPUser("")
	sip.URIType = "sip"
	sip.User = "+12125551212"
	sip.Host = "phone2net.com"
	sip.Args.Set("tag", "887s")
	at.Equal(`<sip:+12125551212@phone2net.com>;tag=887s`, sip.String())

	// Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8
	sip = NewSIPUser("")
	sip.URIType = "sip"
	sip.Name = "Anonymous"
	sip.User = "c8oqz84zk7z"
	sip.Host = "privacy.org"
	sip.Args.Set("tag", "hyh8")
	at.Equal(`Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8`, sip.String())

	// <sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856
	sip = NewSIPUser("")
	sip.URIType = "sip"
	sip.User = "34020000001320000001"
	sip.Host = "192.168.1.102"
	sip.Port = "5060"
	sip.Args.Set("tag", "1086912856")
	at.Equal(`<sip:34020000001320000001@192.168.1.102:5060>;tag=1086912856`, sip.String())
}
