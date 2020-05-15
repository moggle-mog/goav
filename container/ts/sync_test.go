package ts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAvSync(t *testing.T) {
	at := assert.New(t)

	s := newSync(10)
	var dts int64

	dts = 0
	s.syncAudioTs(&dts, 44100)
	at.Equal(int64(0), dts)

	dts = 2000
	s.syncAudioTs(&dts, 44100)
	at.Equal(int64(2089), dts)

	dts = 5000
	s.syncAudioTs(&dts, 44100)
	at.Equal(int64(4178), dts)

	dts = 10000
	s.syncAudioTs(&dts, 44100)
	at.Equal(int64(10000), dts)

	dts = 12000
	s.syncAudioTs(&dts, 44100)
	at.Equal(int64(12089), dts)
}
