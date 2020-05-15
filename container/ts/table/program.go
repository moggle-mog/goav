package table

// Program Ts的节目表
type Program struct {
	Avc []byte
	Aac []byte
}

// NewProgram 新建节目表
func NewProgram() *Program {
	return &Program{
		/*
			stream type: h.264(AVC video stream as defined in ITU-T Rec. H.264 | ISO/IEC 14496-10 Video (h.264))
			stream type pid: 0x100
		*/
		Avc: []byte{0x1b, 0xe1, 0x00, 0xf0, 0x00},
		/*
			stream type: aac(ISO/IEC 13818-7 Audio with ADTS transport syntax)
			stream type pid: 0x101
		*/
		Aac: []byte{0x0f, 0xe1, 0x01, 0xf0, 0x00},
	}
}
