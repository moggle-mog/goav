package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMediaType(t *testing.T) {
	at := assert.New(t)

	mt := NewTypes()
	at.Nil(mt.ToSlice())

	mt.IsVideo()
	at.Equal([]int{PktVideo}, mt.ToSlice())

	mt.Reset()

	mt.IsAudio()
	at.Equal([]int{PktAudio}, mt.ToSlice())

	mt.IsVideo()
	at.Equal([]int{PktVideo, PktAudio}, mt.ToSlice())
}
