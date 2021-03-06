package table

// Pmt Ts的Pmt表
type Pmt struct {
	TsHeader  []byte
	PmtHeader []byte
}

// NewPmt 新建Pmt表
func NewPmt() *Pmt {
	return &Pmt{
		/*
			组成: 4字节固定头 + 1字节指针域
			pid: 0x1001
		*/
		TsHeader: []byte{0x47, 0x50, 0x01, 0x10, 0x00},
		/*
			program number: 0x1
		*/
		PmtHeader: []byte{0x02, 0xb0, 0xff, 0x00, 0x01, 0xc1, 0x00, 0x00, 0xe1, 0x00, 0xf0, 0x00},
	}
}
