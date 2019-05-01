package amf

import (
	"bytes"
	"errors"
	"log"
)

// RTMP辅助方法: 使用AMF协议处理MetaData

// Metadata command
const (
	SetDataFrame string = "@setDataFrame"
	OnMetaData   string = "onMetaData"
)

// setFrameFrame AMF的"SetDataFrame"指令对应的AMF编码
var setFrameFrame []byte

func init() {
	w := bytes.NewBuffer(nil)

	// 将"SetDataFrame"指令所对应的数据封装成AMF0格式, 写到w中
	_, err := NewEnDecAMF0().Encode(w, SetDataFrame)
	if err != nil {
		log.Fatal(err)
	}

	// SetDataFrame对应的AMF0编码
	setFrameFrame = w.Bytes()
}

// AddMetaHeader 在字节序列前加入metadata的头, 用于rtmp传输
func AddMetaHeader(p []byte, d Decoder) ([]byte, error) {
	r := bytes.NewReader(p)

	// 解析AMF编码
	v, err := d.Decode(r)
	if err != nil {
		return nil, err
	}

	// 添加操作, 最后返回[SetDataFrame, p]
	vv, ok := v.(string)
	if !ok {
		return nil, errors.New("setFrameFrame error")
	}

	if vv != SetDataFrame {
		tmpLen := len(setFrameFrame)

		b := make([]byte, tmpLen+len(p))
		copy(b, setFrameFrame)
		copy(b[tmpLen:], p)

		p = b
	}

	return p, nil
}

// AddMetaHeader 从字节序列前删除metadata的头, 用于rtmp传输后解包
func DelMetaHeader(p []byte, d Decoder) ([]byte, error) {
	r := bytes.NewReader(p)

	// 解析AMF编码
	v, err := d.Decode(r)
	if err != nil {
		return nil, err
	}

	// 删除操作, 从[SetDataFrame, p]数据中提取出数据p
	vv, ok := v.(string)
	if !ok {
		return nil, errors.New("metadata error")
	}

	if vv == SetDataFrame {
		p = p[len(setFrameFrame):]
	}

	return p, nil
}
