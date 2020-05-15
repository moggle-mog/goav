package table

// Sdt Ts的Sdt表
type Sdt struct {
	TsHeader  []byte
	SdtHeader []byte
}

// NewSdt 新建Sdt表
func NewSdt() *Sdt {
	return &Sdt{
		/*
			组成: 4字节固定头 + 1字节指针域
			pid: 0x0011
		*/
		TsHeader: []byte{0x47, 0x40, 0x11, 0x10, 0x00},
		/*
			transport stream id: 0x1
		*/
		SdtHeader: []byte{0x42, 0xF0, 0x00, 0x00, 0x01, 0xC1, 0x00, 0x00, 0xFF, 0x01,
			0xFF, 0x00, 0x01, 0xFC, 0x80, 0x00},
	}
}