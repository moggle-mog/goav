package table

// Pat TS的Pat表
type Pat struct {
	TsHeader  []byte
	PatHeader []byte
}

// NewPat 新建Pat表
func NewPat() *Pat {
	return &Pat{
		/*
			组成: 4字节固定头 + 1字节指针域
			pid: 0x0000
		*/
		TsHeader: []byte{0x47, 0x40, 0x00, 0x10, 0x00},
		/*
			program number: 0x1
			program mapping table pid: 0x1001
		*/
		PatHeader: []byte{0x00, 0xb0, 0x0d, 0x00, 0x01, 0xc1, 0x00, 0x00, 0x00, 0x01, 0xf0, 0x01},
	}
}
