package mp3

import (
	"errors"
)

// Parser mp3解析器
type Parser struct {
	samplingFrequency int
}

// NewParser mp3解析器
func NewParser() *Parser {
	return &Parser{
		samplingFrequency: 44100,
	}
}

// sampling_frequency - indicates the sampling frequency, according to the following table.
// '00' 44.1 kHz
// '01' 48 kHz
// '10' 32 kHz
// '11' reserved
var mp3Rates = []int{44100, 48000, 32000}

// Parse 解析mp3数据
func (p *Parser) Parse(src []byte) error {
	if len(src) < 3 {
		return errors.New("incomplete mp3 data, len(src)<3")
	}

	// 提取出采样率
	index := (src[2] >> 1) & 0x3
	if int(index) < len(mp3Rates) {
		p.samplingFrequency = mp3Rates[index]
		return nil
	}

	return errors.New("invalid rate index")
}

// SampleRate mp3采样率
func (p *Parser) SampleRate() int {
	return p.samplingFrequency
}
