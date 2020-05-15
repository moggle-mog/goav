package flv

import (
	"bytes"
	"testing"

	"github.com/moggle-mog/goav/packet"

	"github.com/moggle-mog/goav/amf"
	"github.com/stretchr/testify/assert"
)

func TestMuxer_Mux(t *testing.T) {
	at := assert.New(t)

	mux := newMuxer()

	buf := bytes.NewBuffer(nil)

	data, err := mux.header(packet.PktVideo)
	at.Nil(err)
	buf.Write(data)

	data, err = mux.metadata(amf.Object{
		"Provider": "test provider",
	})
	at.Nil(err)
	buf.Write(data)

	at.Equal([]byte{
		0x46, 0x4c, 0x56, 0x1, 0x1, 0x0, 0x0, 0x0,
		0x9, 0x0, 0x0, 0x0, 0x0, 0x12, 0x0, 0x0,
		0x2b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x2, 0x0, 0xa, 0x6f, 0x6e, 0x4d, 0x65, 0x74,
		0x61, 0x44, 0x61, 0x74, 0x61, 0x3, 0x0, 0x8,
		0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
		0x2, 0x0, 0xd, 0x74, 0x65, 0x73, 0x74, 0x20,
		0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
		0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x36,
	}, buf.Bytes())
}