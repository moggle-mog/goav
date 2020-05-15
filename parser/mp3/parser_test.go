package mp3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_SampleRate(t *testing.T) {
	at := assert.New(t)

	p := NewParser()
	at.Nil(p.Parse([]byte{0x62, 0x70, 0x6C}))

	at.Equal(32000, p.SampleRate())
}
