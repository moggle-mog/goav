package ts

import (
	"bytes"
	"testing"

	"github.com/moggle-mog/goav/container/ts/table"
	"github.com/moggle-mog/goav/packet"
	"github.com/stretchr/testify/assert"
)

type TestWriter struct {
	buf   []byte
	count int
}

func (w *TestWriter) Write(p []byte) (int, error) {
	w.count++
	w.buf = p
	return len(p), nil
}

func TestTSEncoder(t *testing.T) {
	at := assert.New(t)

	m := NewMuxer()

	w := &TestWriter{}
	media := []byte{
		0xaf, 0x01, 0x21, 0x19, 0xd3, 0x40, 0x7d, 0x0b,
		0x6d, 0x44, 0xae, 0x81, 0x08, 0x00, 0x89, 0xa0,
		0x3e, 0x85, 0xb6, 0x92, 0x57, 0x04, 0x80, 0x00,
		0x5b, 0xb7, 0x78, 0x00, 0x84, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x38, 0x30, 0x00, 0x06, 0x00, 0x38,
	}

	p := &packet.Packet{
		Type:  packet.PktAudio,
		Media: media,
	}

	err := m.Mux(p, 0, 0, w)
	at.Equal(nil, err)
	at.Equal(1, w.count)

	expected := []byte{
		0x47, 0x41, 0x01, 0x31, 0x81, 0x00, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00,
		0x01, 0xc0, 0x00, 0x30, 0x80, 0x80, 0x05, 0x21,
		0x00, 0x01, 0x00, 0x01, 0xaf, 0x01, 0x21, 0x19,
		0xd3, 0x40, 0x7d, 0x0b, 0x6d, 0x44, 0xae, 0x81,
		0x08, 0x00, 0x89, 0xa0, 0x3e, 0x85, 0xb6, 0x92,
		0x57, 0x04, 0x80, 0x00, 0x5b, 0xb7, 0x78, 0x00,
		0x84, 0x00, 0x00, 0x00, 0x00, 0x00, 0x38, 0x30,
		0x00, 0x06, 0x00, 0x38,
	}
	at.Equal(expected, w.buf)
}

func TestMuxer_Table(t *testing.T) {
	at := assert.New(t)

	mux := NewMuxer()
	buf := bytes.NewBuffer(nil)

	sdtDesc := table.NewDescriptor()
	at.Nil(sdtDesc.Service(1, "test provider", "test service"))

	// 输出SDT表
	sdt := mux.SDT(sdtDesc.GetBuffer())
	at.NotNil(sdt)

	_, err := buf.Write(sdt)
	at.Nil(err)

	// 输出PAT表
	pat := mux.PAT()
	at.NotNil(pat)

	_, err = buf.Write(pat)
	at.Nil(err)

	// 输出PMT表
	pmt := mux.PMT(packet.PktVideo)
	_, err = buf.Write(pmt)
	at.Nil(err)

	at.Equal([]byte{
		0x47, 0x40, 0x11, 0x10, 0x0, 0x42, 0xf0, 0x2f,
		0x0, 0x1, 0xc1, 0x0, 0x0, 0xff, 0x1, 0xff, 0x0,
		0x1, 0xfc, 0x80, 0x1e, 0x48, 0x1c, 0x1, 0xd,
		0x74, 0x65, 0x73, 0x74, 0x20, 0x70, 0x72, 0x6f,
		0x76, 0x69, 0x64, 0x65, 0x72, 0xc, 0x74, 0x65,
		0x73, 0x74, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69,
		0x63, 0x65, 0xf0, 0x3d, 0x5a, 0xe2, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0x47, 0x40, 0x0, 0x10, 0x0,
		0x0, 0xb0, 0xd, 0x0, 0x1, 0xc1, 0x0, 0x0,
		0x0, 0x1, 0xf0, 0x1, 0x2e, 0x70, 0x19, 0x5,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x47,
		0x50, 0x1, 0x10, 0x0, 0x2, 0xb0, 0x12, 0x0,
		0x1, 0xc1, 0x0, 0x0, 0xe1, 0x0, 0xf0, 0x0,
		0x1b, 0xe1, 0x0, 0xf0, 0x0, 0x15, 0xbd, 0x4d,
		0x56, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff,
	}, buf.Bytes())
}
