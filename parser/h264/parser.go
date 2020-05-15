package h264

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Parser H264解析器
type Parser struct {
	specificInfo []byte        /* {0: sps, 1: pps}, 均包含start code */
	spsPps       *bytes.Buffer /* sps和pps共用, 均包含start code */
}

// NewParser 初始化h264解析器(pps/sps)
func NewParser() *Parser {
	return &Parser{
		spsPps: bytes.NewBuffer(make([]byte, maxSpsPpsLen)),
	}
}

// Parse 将H264打包格式转换为 Annex-b 的网络流格式, 写入w中
func (p *Parser) Parse(b []byte, isSeqHdr bool, w io.Writer) error {
	if len(b) == 0 || w == nil {
		return errors.New("no data to parse or nil writer")
	}

	// [AVCC格式]如果是序列头, 则解析出SPS和PPS
	if isSeqHdr {
		return p.parseSpecificInfo(b)
	}

	// [Annex-b格式]直接写入以Nalu开头的数据
	if p.isStartAtNaluHeader(b) {
		_, err := w.Write(b)
		if err != nil {
			return err
		}

		return nil
	}

	// [AVCC格式]转换为Annex-b格式并写入数据
	return p.getAnnexbH264(b, w)
}

// [AVCC格式]向specificInfo填充SPS和PPS, specificInfo的值: {0: sps数据, 1: pps数据}
func (p *Parser) parseSpecificInfo(src []byte) error {
	if len(src) < 9 {
		return errors.New("incomplete data, len(src)<9")
	}

	var sps []byte
	var pps []byte

	var seq sequenceHeader

	// 填充 AVCDecoderConfigurationRecord
	seq.configurationVersion = src[0] /* NALU头, 版本号: 1 */
	seq.avcProfileIndication = src[1] /* H264的profile, Baseline: 视频会议, Main: 标清电视, High: 高清电视 */
	seq.profileCompatility = src[2]
	seq.avcLevelIndication = src[3] /* h264的Level, 可以指定最大分辨率, 帧率等 */
	seq.reserved1 = src[4] >> 2
	seq.naluLen = src[4]&0x3 + 1
	seq.reserved2 = src[5] >> 5

	// 提取SPS
	seq.spsNum = src[5] & 0x1f                /* [3:7]SPS数量, 一般为1 */
	seq.spsLen = int(src[6])<<8 | int(src[7]) /* [0:15]SPS长度 */
	if len(src[8:]) < seq.spsLen || seq.spsLen <= 0 {
		return errors.New("incomplete sps data")
	}
	sps = append(sps, startCode...)
	sps = append(sps, src[8:(8+seq.spsLen)]...)

	// 提取PPS
	tmpBuf := src[(8 + seq.spsLen):]
	if len(tmpBuf) < 4 {
		return errors.New("incomplete pps header")
	}
	seq.ppsNum = tmpBuf[0]                          /* PPS数量 */
	seq.ppsLen = int(tmpBuf[1])<<8 | int(tmpBuf[2]) /* [0:15]PPS长度 */
	if len(tmpBuf[3:]) < seq.ppsLen || seq.ppsLen <= 0 {
		return errors.New("incomplete pps data")
	}
	pps = append(pps, startCode...)
	pps = append(pps, tmpBuf[3:]...)

	// 向specificInfo填充SPS和PPS
	p.specificInfo = append(p.specificInfo, sps...)
	p.specificInfo = append(p.specificInfo, pps...)

	return nil
}

// 判断数据是否是以NALU头开始, Annex-b格式以NALU头开始
func (p *Parser) isStartAtNaluHeader(src []byte) bool {
	if len(src) < naluBytesLen {
		return false
	}

	return src[0] == 0x00 && src[1] == 0x00 && src[2] == 0x00 && src[3] == 0x01
}

// [AVCC格式] 提取NALU的长度
func (p *Parser) naluSize(src []byte) (int, error) {
	if len(src) < naluBytesLen {
		return 0, errors.New("[hvcc]incomplete nalu data")
	}

	// [0:31] nalu单元的长度，不包括长度字段
	buf := src[:naluBytesLen]
	l := len(buf)

	var size = 0
	for i := 0; i < l; i++ {
		size = size<<8 + int(buf[i])
	}

	return size, nil
}

// [AVCC->Annex-b]将以 AVCC 作为打包格式转换为以 Annex-b 作为打包格式的H264数据写入w中
func (p *Parser) getAnnexbH264(src []byte, w io.Writer) error {
	dataSize := len(src)

	if dataSize == 0 || dataSize < naluBytesLen {
		return errors.New("incomplete h264 header")
	}

	// 写入音频的nalu
	_, err := w.Write(naluAud)
	if err != nil {
		return err
	}

	index := 0
	hasSpsPps := false
	hasWriteSpsPps := false

	// 重置sps/pps的值
	p.spsPps.Reset()

	// 从AVCC的打包格式转换为Annex-b的打包格式; 对于整个流程, 首先写入SPS和PPS, 然后都只是向后填充数据
	for dataSize > 0 {
		// 	取出nalu的size
		nalLen, err := p.naluSize(src[index:])
		if err != nil {
			return err
		}

		if nalLen <= 0 {
			return errors.New("invalid nalu body size")
		}

		index += naluBytesLen
		dataSize -= naluBytesLen

		// 紧跟着 Nalu size 后面的是 NALU 数据，没有四个字节 start code，直接从 h264 头开始
		if dataSize < nalLen {
			return errors.New("invalid nalu body")
		}

		nalType := src[index] & 0x1f /* [3:7]nal_unit_type 帧类型 */

		switch nalType {
		case naluTypeAud:
		case naluTypeIdr:
			// 如果未接入SPS和PPS信息, 则写入SPS和PPS,
			// 如果视频包中有SPS或者PPS, 则从视频包提取该数据, 否则从缓存的数据里提取SPS和PPS
			if !hasWriteSpsPps {
				hasWriteSpsPps = true
				if hasSpsPps {
					// 根据pps, 写入 SPS/PPS
					_, err = w.Write(p.spsPps.Bytes())
					if err != nil {
						return err
					}
				} else {
					// 根据specificInfo, 写入 SPS和PPS
					_, err = w.Write(p.specificInfo)
					if err != nil {
						return err
					}
				}
			}
			fallthrough
		case naluTypeSlice:
			fallthrough
		case naluTypeSei:
			// 写入 start code
			_, err = w.Write(startCode)
			if err != nil {
				return err
			}

			// 写入 sei 数据
			_, err = w.Write(src[index : index+nalLen])
			if err != nil {
				return err
			}
		case naluTypeSps:
			fallthrough
		case naluTypePps:
			hasSpsPps = true

			// 写入 start code
			_, err = p.spsPps.Write(startCode)
			if err != nil {
				return err
			}

			// 写入 SPS 或 PPS 的数据
			_, err = p.spsPps.Write(src[index : index+nalLen])
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("incompatible nalu type number=%d", nalType)
		}

		index += nalLen
		dataSize -= nalLen
	}

	return nil
}
