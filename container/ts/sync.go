package ts

// 音视频频率
const (
	// AACSL AAC的频率
	aacSL = 1024

	// AVCHZ H264的频率
	avcHZ = 90
)

// sync 音视频同步
type sync struct {
	frameNum int64 // 累计帧数
	frameDts int64 // 基准时间
	syncMs   int64 // ms, 同步 |pts-dts|>syncMs 的"pts"和"dts"
}

// newSync 音视频时间戳同步
func newSync(ms int64) *sync {
	if ms <= 0 {
		panic("ms<=0")
	}
	return &sync{
		syncMs: ms * avcHZ,
	}
}

// SyncAudioTs 音视频同步，根据视频dts时间调整音频时间
// dts: 传入音频的解码时间戳, 传出音频的播放时间戳, 单位: ms
// sampleRate: 音频采样率, 单位: HZ
func (s *sync) syncAudioTs(dts *int64, sampleRate int) {
	// 根据采样率, 换算音频相对于视频的时间增量
	tsIncrement := avcHZ * 1000 * aacSL / sampleRate
	pts := s.frameDts + s.frameNum*int64(tsIncrement)

	// 计算出pts和dts之间的差值
	var ptsDtsGap int64
	if pts >= *dts {
		ptsDtsGap = pts - *dts
	} else {
		ptsDtsGap = *dts - pts
	}

	// 差值在阈值内，dts=pts
	if ptsDtsGap <= s.syncMs {
		s.frameNum++
		*dts = pts
		return
	}

	// 差值在阈值外，dts=dts
	s.frameNum = 1
	s.frameDts = *dts
}
