package packet

// FLV Tag
const (
	TagAudio          = 8   // Tag：视频
	TagVideo          = 9   // Tag：音频
	TagScriptDataAMF0 = 18  // Tag：AMF0格式的元数据
	TagScriptDataAMF3 = 0xf // Tag：AMF3格式的元数据
)

// FLV package type
const (
	PktVideo    = iota // 视频包
	PktAudio           // 音频包
	PktMetadata        // 元数据包
)

// Packet Header can be converted to AudioHeaderInfo or VideoHeaderInfo
type Packet struct {
	Type      int    // 音频, 视频, 元数据, 其它
	TimeStamp uint32 // dts, 增量, 毫秒
	Baseline  uint32 // 增量，毫秒
	StreamID  uint32 // 流ID
	Header    Header // FLV头
	Data      []byte // 封装了flv tag的数据包
	Media     []byte // 裸流数据
}

// PacketHeader 帧头
type Header interface{}

// AudioPacketHeader FLV音频帧描述接口
type AudioPacketHeader interface {
	Header
	SoundFormat() uint8
	AACType() uint8
	IsSoundAAC() bool
	IsSoundMP3() bool
	IsAACSeqHdr() bool
}

// VideoPacketHeader FLV视频帧描述接口
type VideoPacketHeader interface {
	Header
	IsKeyFrame() bool
	IsInterFrame() bool
	IsSeqHdr() bool
	IsEndOfSeq() bool
	IsCodecAvc() bool
	CodecID() uint8
	CompositionTime() int32
}
