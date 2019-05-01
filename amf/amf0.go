package amf

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"
)

// AMF0 Data Type(U8): https://wwwimages2.adobe.com/content/dam/acom/en/devnet/pdf/amf0-file-format-specification.pdf
const (
	numberMarkerOnAMF0        = 0x00
	booleanMarkerOnAMF0       = 0x01
	stringMarkerOnAMF0        = 0x02
	objectMarkerOnAMF0        = 0x03
	movieclipMarkerOnAMF0     = 0x04 /* reserved, not used */
	nullMarkerOnAMF0          = 0x05
	undefinedMarkerOnAMF0     = 0x06
	referenceMarkerOnAMF0     = 0x07
	ecmaArrayMarkerOnAMF0     = 0x08
	objectEndMarkerOnAMF0     = 0x09
	strictArrayMarkerOnAMF0   = 0x0a
	dateMarkerOnAMF0          = 0x0b
	longStringMarkerOnAMF0    = 0x0c
	unsupportedMarkerOnAMF0   = 0x0d
	recordsetMarkerOnAMF0     = 0x0e /* reserved, not used */
	xmlDocumentMarkerOnAMF0   = 0x0f
	typedObjectMarkerOnAMF0   = 0x10
	avmPlusObjectMarkerOnAMF0 = 0x11 /* switch to AMF3 */
)

// EnDecAMF0 AMF0协议解析器
type EnDecAMF0 struct {
	deRefCache []interface{}
	EnDecAMF3
}

// NewEnDecAMF0 AMF0协议解析器
func NewEnDecAMF0() *EnDecAMF0 {
	return &EnDecAMF0{}
}

// Decode AMF0协议解析器入口
func (ed *EnDecAMF0) Decode(r io.Reader) (interface{}, error) {
	// 读取第一个字节(包类型)
	marker, err := readByte(r)
	if err != nil {
		return nil, err
	}

	switch marker {
	case numberMarkerOnAMF0:
		// marker: 1 byte 0x00
		// format: 8 byte big endian float64
		var f64 float64
		err = binary.Read(r, binary.BigEndian, &f64)
		if err != nil {
			return f64, err
		}

		return f64, nil
	case booleanMarkerOnAMF0:
		// marker: 1 byte 0x01
		// format: 1 byte, 0x00 = false, 0x01 = true
		var flag byte
		flag, err = readByte(r)
		if err != nil {
			return nil, err
		}

		return flag == 1, nil
	case stringMarkerOnAMF0:
		return ed.decodeString(r)
	case objectMarkerOnAMF0:
		return ed.decodeObject(r)
	case movieclipMarkerOnAMF0:
		// marker: 1 byte 0x04
		// This type is not supported and is reserved for future use
		return nil, fmt.Errorf("decode amf0: unsupported type movieclip")
	case nullMarkerOnAMF0, undefinedMarkerOnAMF0, unsupportedMarkerOnAMF0:
		return nil, nil
	case referenceMarkerOnAMF0:
		return ed.decodeReference(r)
	case ecmaArrayMarkerOnAMF0:
		return ed.decodeEcmaArray(r)
	case strictArrayMarkerOnAMF0:
		return ed.decodeStrictArray(r)
	case dateMarkerOnAMF0:
		return ed.decodeDate(r)
	case longStringMarkerOnAMF0:
		return ed.decodeLongString(r)
	case recordsetMarkerOnAMF0:
		return nil, fmt.Errorf("decode amf0: unsupported type recordset")
	case xmlDocumentMarkerOnAMF0:
		// marker: 1 byte 0x0f
		// format:
		// - normal long string format
		//   - 4 byte big endian uint32 header to determine size
		//   - n (size) byte utf8 string
		return ed.decodeLongString(r)
	case typedObjectMarkerOnAMF0:
		return ed.decodeTypedObject(r)
	case avmPlusObjectMarkerOnAMF0:
		return ed.EnDecAMF3.Decode(r)
	}

	return nil, fmt.Errorf("decode amf0: unsupported type %d", marker)
}

// DecodeBatch 批量解析AMF(底层调用Decode)
func (ed *EnDecAMF0) DecodeBatch(r io.Reader) ([]interface{}, error) {
	var ret []interface{}
	for {
		v, err := ed.Decode(r)
		if err != nil {
			if err == io.EOF {
				return ret, nil
			}

			return ret, err
		}

		ret = append(ret, v)
	}
}

// ====================================================================

func (ed *EnDecAMF0) decodeString(r io.Reader) (string, error) {
	// 字符串最大长度: 65535(B) ≈ 64(K)
	var length uint16
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}

	// 协议要求读取的字节数一定要满足, 否则输出错误, 至于超过的可以忽略
	buf, err := readBytes(r, int(length))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (ed *EnDecAMF0) decodeObject(r io.Reader) (Object, error) {
	result := make(Object)
	ed.deRefCache = append(ed.deRefCache, result)

	for {
		// 读取字符串
		key, err := ed.decodeString(r)
		if err != nil {
			return nil, err
		}

		// [出口]空字符串后应该是结束符
		if key == "" {
			err = nextByteMustBe(r, objectEndMarkerOnAMF0)
			if err != nil {
				return nil, err
			}

			return result, nil
		}

		// 序列未结束, 继续读取
		value, err := ed.Decode(r)
		if err != nil {
			return nil, err
		}

		result[key] = value
	}
}

func (ed *EnDecAMF0) decodeReference(r io.Reader) (interface{}, error) {
	var ref uint16
	err := binary.Read(r, binary.BigEndian, &ref)
	if err != nil {
		return nil, err
	}

	if int(ref) > len(ed.deRefCache) {
		return nil, fmt.Errorf("bad reference %d (current length %d)", ref, len(ed.deRefCache))
	}

	return ed.deRefCache[ref], nil
}

func (ed *EnDecAMF0) decodeEcmaArray(r io.Reader) (Object, error) {
	// 关联数组长度
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	// 解析对象
	result, err := ed.decodeObject(r)
	if err != nil {
		return nil, err
	}

	// length作为校验手段未使用上

	return result, nil
}

func (ed *EnDecAMF0) decodeStrictArray(r io.Reader) (Array, error) {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	// 初始化数组
	var result Array
	ed.deRefCache = append(ed.deRefCache, result)

	// 填充数组
	for i := uint32(0); i < length; i++ {
		val, err := ed.Decode(r)
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

func (ed *EnDecAMF0) decodeDate(r io.Reader) (float64, error) {
	var result float64
	err := binary.Read(r, binary.BigEndian, &result)
	if err != nil {
		return 0, err
	}

	// 读取2字节(未使用)
	_, err = readBytes(r, 2)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (ed *EnDecAMF0) decodeLongString(r io.Reader) (string, error) {
	// 字符串最大长度: 4294967295 Bytes ≈ 4G
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}

	buf, err := readBytes(r, int(length))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (ed *EnDecAMF0) decodeTypedObject(r io.Reader) (TypedObject, error) {
	// 数组化类型数组
	var result TypedObject
	ed.deRefCache = append(ed.deRefCache, result)

	// 读取类型
	typeVal, err := ed.decodeString(r)
	if err != nil {
		return result, err
	}

	// 读取对象
	objectVal, err := ed.decodeObject(r)
	if err != nil {
		return result, err
	}

	// 导出返回值
	result = TypedObject{
		Type:   typeVal,
		Object: objectVal,
	}

	return result, nil
}

// ====================================================================

// Encode AMF0协议编码器入口
// 支持"Object", "TypedObject", "time", "nil", "float", "bool", "string", "int", "uint", "array", "slice"
// 其他类型请直接调用对应的编码器
func (ed *EnDecAMF0) Encode(w io.Writer, val interface{}) (int, error) {
	// 取得变量信息
	v := reflect.ValueOf(val)

	// 复杂类型
	switch val.(type) {
	case Object:
		return ed.encodeObject(w, val.(Object), true)
	case TypedObject:
		return ed.encodeTypeObject(w, val.(TypedObject))
	case time.Time:
		return ed.encodeDate(w, val.(time.Time))
	}

	// 简单类型
	switch v.Kind() {
	case reflect.Invalid:
		return ed.encodeNull(w)
	case reflect.Float32, reflect.Float64:
		return ed.encodeNumber(w, float64(v.Float()))
	case reflect.Bool:
		return ed.encodeBoolean(w, v.Bool())
	case reflect.String:
		if v.Len() <= 65535 {
			return ed.encodeString(w, v.String(), true)
		}

		return ed.encodeLongString(w, v.String(), true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ed.encodeNumber(w, float64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ed.encodeNumber(w, float64(v.Uint()))
	case reflect.Array, reflect.Slice:
		arr := make(Array, v.Len())
		for i := 0; i < v.Len(); i++ {
			arr[i] = v.Index(int(i)).Interface()
		}

		return ed.encodeStrictArray(w, arr)
	}

	// 未知类型
	return 0, fmt.Errorf("undefined type %s", v.Type())
}

// EncodeBatch 批量封装AMF(底层调用Encode)
func (ed *EnDecAMF0) EncodeBatch(w io.Writer, args ...interface{}) error {
	var err error
	for _, v := range args {
		_, err = ed.Encode(w, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// EncodeWithAMF3 使用AMF3编码数据
func (ed *EnDecAMF0) EncodeWithAMF3(w io.Writer, val interface{}) (int, error) {
	err := writeByte(w, avmPlusObjectMarkerOnAMF0)
	if err != nil {
		return 0, nil
	}

	n, err := ed.EnDecAMF3.Encode(w, val)
	if err != nil {
		return 0, err
	}

	return n + 1, nil
}

// EncodeXmlDocument 编码Xml类型
func (ed *EnDecAMF0) EncodeXmlDocument(w io.Writer, val string) (int, error) {
	err := writeByte(w, xmlDocumentMarkerOnAMF0)
	if err != nil {
		return 0, nil
	}

	n, err := ed.encodeLongString(w, val, false)
	if err != nil {
		return 0, err
	}

	return n + 1, nil
}

// EncodeReference 编写引用类型
func (ed *EnDecAMF0) EncodeReference(w io.Writer, ref uint16) (int, error) {
	err := writeByte(w, referenceMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	err = binary.Write(w, binary.BigEndian, ref)
	if err != nil {
		return 0, err
	}
	n += 2

	return n, nil
}

// EncodeEcmaArray 编码关联数组类型
func (ed *EnDecAMF0) EncodeEcmaArray(w io.Writer, val Object) (int, error) {
	err := writeByte(w, ecmaArrayMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	// 数组长度
	length := uint32(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	n += 4

	// 数组(对象)
	m, err := ed.encodeObject(w, val, false)
	if err != nil {
		return 0, err
	}
	n += m

	return n, nil
}

// EncodeUndefined 编码未定义类型
func (ed *EnDecAMF0) EncodeUndefined(w io.Writer) (int, error) {
	err := writeByte(w, undefinedMarkerOnAMF0)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// EncodeUndefined 编码不支持类型
func (ed *EnDecAMF0) EncodeUnsupported(w io.Writer) (int, error) {
	err := writeByte(w, unsupportedMarkerOnAMF0)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// ====================================================================

func (ed *EnDecAMF0) encodeNumber(w io.Writer, val float64) (int, error) {
	err := writeByte(w, numberMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	err = binary.Write(w, binary.BigEndian, &val)
	if err != nil {
		return 0, err
	}
	n += 8

	return n, nil
}

func (ed *EnDecAMF0) encodeBoolean(w io.Writer, val bool) (int, error) {
	err := writeByte(w, booleanMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	buf := make([]byte, 1)
	if val {
		buf[0] = 1
	} else {
		buf[0] = 0
	}

	_, err = w.Write(buf)
	if err != nil {
		return 0, err
	}
	n += 1

	return n, nil
}

func (ed *EnDecAMF0) encodeString(w io.Writer, val string, encodeMarker bool) (int, error) {
	var n int
	var err error

	if encodeMarker {
		err = writeByte(w, stringMarkerOnAMF0)
		if err != nil {
			return 0, err
		}
		n++
	}

	length := uint16(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	n += 2

	_, err = w.Write([]byte(val))
	if err != nil {
		return 0, err
	}
	n += int(length)

	return n, nil
}

func (ed *EnDecAMF0) encodeObject(w io.Writer, val Object, encodeMarker bool) (int, error) {
	var n int
	var err error

	if encodeMarker {
		err = writeByte(w, objectMarkerOnAMF0)
		if err != nil {
			return 0, err
		}
		n++
	}

	var m = 0
	for k, v := range val {
		// 数组key
		m, err = ed.encodeString(w, k, false)
		if err != nil {
			return 0, err
		}
		n += m

		// 数组value
		m, err = ed.Encode(w, v)
		if err != nil {
			return 0, err
		}
		n += m
	}

	// 结束符: UTF-8-empty
	m, err = ed.encodeString(w, "", false)
	if err != nil {
		return 0, err
	}
	n += m

	// 结束符: object-end-marker
	err = writeByte(w, objectEndMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n++

	return n, nil
}

func (ed *EnDecAMF0) encodeNull(w io.Writer) (int, error) {
	err := writeByte(w, nullMarkerOnAMF0)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (ed *EnDecAMF0) encodeStrictArray(w io.Writer, val Array) (int, error) {
	err := writeByte(w, strictArrayMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	// 数组长度
	length := uint32(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	n += 4

	// 数组(多类型)
	var m = 0
	for _, v := range val {
		m, err = ed.Encode(w, v)
		if err != nil {
			return 0, err
		}
		n += m
	}

	return n, nil
}

func (ed *EnDecAMF0) encodeDate(w io.Writer, val time.Time) (int, error) {
	err := writeByte(w, dateMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	var m = 0
	m, err = ed.encodeNumber(w, float64(val.Unix())*1000.0)
	if err != nil {
		return 0, err
	}
	n += m

	err = binary.Write(w, binary.BigEndian, int16(0))
	if err != nil {
		return 0, err
	}
	n += 2

	return n, nil
}

func (ed *EnDecAMF0) encodeLongString(w io.Writer, val string, encodeMarker bool) (int, error) {
	var n int
	var err error

	if encodeMarker {
		err = writeByte(w, longStringMarkerOnAMF0)
		if err != nil {
			return 0, err
		}
		n++
	}

	// 写入字符串长度
	length := uint32(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	n += 4

	// 写入字符串
	m, err := w.Write([]byte(val))
	if err != nil {
		return 0, err
	}
	n += m

	return n, nil
}

func (ed *EnDecAMF0) encodeTypeObject(w io.Writer, tyeObj TypedObject) (int, error) {
	err := writeByte(w, typedObjectMarkerOnAMF0)
	if err != nil {
		return 0, err
	}
	n := 1

	var m = 0
	m, err = ed.encodeString(w, tyeObj.Type, false)
	if err != nil {
		return 0, err
	}
	n += m

	m, err = ed.encodeObject(w, tyeObj.Object, false)
	if err != nil {
		return 0, err
	}
	n += m

	return n, nil
}
