package aac

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAac(t *testing.T) {
	at := assert.New(t)
	d := NewParser()
	w := bytes.NewBuffer(nil)

	err := d.Parse([]byte{0x12, 0x10}, SeqHdr, w)
	at.Equal(nil, err)

	audio := []byte{
		0x21, 0x00, 0x49, 0x90, 0x02, 0x19, 0x00, 0x23, 0x80,
	}
	err = d.Parse(audio, Raw, w)
	at.Equal(nil, err)
	at.Equal([]byte{0xff, 0xf1, 0x50, 0x80, 0x02, 0x1f, 0xfc, 0x21, 0x00, 0x49, 0x90, 0x02, 0x19, 0x00, 0x23, 0x80}, w.Bytes())
}