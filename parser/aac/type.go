package aac

const adtsHeaderLen = 7

// AAC type
const (
	SeqHdr = iota
	Raw
)

var aacRates = []int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350}

type mpegExtension struct {
	objectType      byte
	sampleRateIndex byte
}

type mpegCfgInfo struct {
	objectType      byte // aac sequence header: 编码结构类型
	sampleRateIndex byte // aac sequence header: 音频采样率索引值
	channel         byte // aac sequence header: 音频输出声道
	sbr             byte
	ps              byte
	frameLen        byte
	exceptionLogTs  int64
	extension       *mpegExtension
}
