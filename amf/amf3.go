package amf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"time"
)

// AMF3 Data Type(U8): https://wwwimages2.adobe.com/content/dam/acom/en/devnet/pdf/amf-file-format-spec.pdf
const (
	undefinedMarkerOnAMF3    = 0x00
	nullMarkerOnAMF3         = 0x01
	falseMarkerOnAMF3        = 0x02
	trueMarkerOnAMF3         = 0x03
	integerMarkerOnAMF3      = 0x04
	doubleMarkerOnAMF3       = 0x05
	stringMarkerOnAMF3       = 0x06
	xmlDocMarkerOnAMF3       = 0x07
	dateMarkerOnAMF3         = 0x08
	arrayMarkerOnAMF3        = 0x09
	objectMarkerOnAMF3       = 0x0a
	xmlMarkerOnAMF3          = 0x0b
	byteArrayMarkerOnAMF3    = 0x0c
	vectorIntMarkerOnAMF3    = 0x0d
	vectorUintMarkerOnAMF3   = 0x0e
	vectorDoubleMarkerOnAMF3 = 0x0f
	vectorObjectMarkerOnAMF3 = 0x10
	dictionaryMarkerOnAMF3   = 0x11
)

// TraitAMF3 AMF3特性表
type TraitAMF3 struct {
	Type           string
	Externalizable bool
	Dynamic        bool
	Properties     []string
}

// EcmaArray AMF3关联数组接口
type EcmaArrayAMF3 struct {
	dense       Array
	associative Object
}

// ExternalHandler 远程调用
type ExternalHandler func(*EnDecAMF3, io.Reader) (interface{}, error)

// EnDecAMF3 AMF3协议解析器
type EnDecAMF3 struct {
	deStrTable       []string
	deObjTable       []interface{}
	traitTable       []TraitAMF3
	externalHandlers map[string]ExternalHandler
}

// NewEnDecAMF3 AMF3协议解析器
func NewEnDecAMF3() *EnDecAMF3 {
	return &EnDecAMF3{
		externalHandlers: make(map[string]ExternalHandler),
	}
}

// RegisterExternalHandler 注册远程调用
func (ed *EnDecAMF3) RegisterExternalHandler(name string, f ExternalHandler) {
	ed.externalHandlers[name] = f
}

// Decode AMF3解析器入口
func (ed *EnDecAMF3) Decode(r io.Reader) (interface{}, error) {
	marker, err := readByte(r)
	if err != nil {
		return nil, err
	}

	switch marker {
	case undefinedMarkerOnAMF3, nullMarkerOnAMF3:
		return nil, nil
	case falseMarkerOnAMF3:
		return false, nil
	case trueMarkerOnAMF3:
		return true, nil
	case integerMarkerOnAMF3:
		return ed.readUint29(r)
	case doubleMarkerOnAMF3:
		var f64 float64
		if err := binary.Read(r, binary.BigEndian, &f64); err != nil {
			return f64, err
		}

		return f64, nil
	case stringMarkerOnAMF3:
		return ed.decodeString(r)
	case xmlDocMarkerOnAMF3:
		return ed.decodeXmlDoc(r)
	case dateMarkerOnAMF3:
		return ed.decodeDate(r)
	case arrayMarkerOnAMF3:
		return ed.decodeArray(r)
	case objectMarkerOnAMF3:
		return ed.decodeObject(r)
	case xmlMarkerOnAMF3:
		return ed.decodeXmlDoc(r)
	case byteArrayMarkerOnAMF3:
		return ed.decodeByteArray(r)
	case vectorIntMarkerOnAMF3, vectorUintMarkerOnAMF3, vectorDoubleMarkerOnAMF3, vectorObjectMarkerOnAMF3:
		return nil, fmt.Errorf("decode amf3: unrealized type %d", marker)
	case dictionaryMarkerOnAMF3:
		return nil, fmt.Errorf("decode amf3: unrealized type %d", marker)
	}

	return nil, fmt.Errorf("decode amf3: unsupported type %d", marker)
}

// DecodeBatch 批量解析AMF(底层调用Decode)
func (ed *EnDecAMF3) DecodeBatch(r io.Reader) ([]interface{}, error) {
	var ret []interface{}
	for {
		v, err := ed.Decode(r)
		if err != nil {
			return ret, err
		}
		ret = append(ret, v)
	}
}

// ======================================================================

func (ed *EnDecAMF3) decodeString(r io.Reader) (string, error) {
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return "", err
	}

	// If the flag is 0, then a string reference is encoded and
	// the remaining bits are used to encode an index to the implicit string reference table
	if isRef {
		index := int(refVal)
		if index >= len(ed.deStrTable) {
			return "", fmt.Errorf("amf3 decode: string index outbound(expected < %d, got %d)", len(ed.deStrTable), index)
		}

		return ed.deStrTable[index], nil
	}

	// 待读字节数
	if refVal == 0 {
		return "", nil
	}

	// If the flag is 1, a string literal is encoded and
	// the remaining bits are used to encode the byte-length of the UTF-8 encoded String
	buf, err := readBytes(r, int(refVal))
	if err != nil {
		return "", err
	}

	// 写入缓存
	result := string(buf)
	ed.deStrTable = append(ed.deStrTable, result)

	return result, nil
}

func (ed *EnDecAMF3) decodeDate(r io.Reader) (time.Time, error) {
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return time.Time{}, err
	}

	// 按引用取值
	if isRef {
		index := int(refVal)
		if index >= len(ed.deObjTable) {
			return time.Time{}, fmt.Errorf("amf3 decode date: object index outbound(expected < %d, got %d)", len(ed.deObjTable), index)
		}

		// 将数值转换为time类型
		result, ok := ed.deObjTable[index].(time.Time)
		if !ok {
			return time.Time{}, fmt.Errorf("amf3 decode date: unable to extract time from date object references")
		}

		return result, nil
	}

	// 按数值取值
	var f64 float64
	if err = binary.Read(r, binary.BigEndian, &f64); err != nil {
		return time.Time{}, fmt.Errorf("amf3 decode date: unable to read double: %s", err)
	}

	// 将毫秒转换为秒
	result := time.Unix(int64(f64/1000), 0).UTC()
	ed.deObjTable = append(ed.deObjTable, result)

	return result, nil
}

func (ed *EnDecAMF3) decodeArray(r io.Reader) (EcmaArrayAMF3, error) {
	// 解析命令
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return EcmaArrayAMF3{}, err
	}

	// 按引用取值(deObjTable)
	if isRef {
		index := int(refVal)
		if index >= len(ed.deObjTable) {
			return EcmaArrayAMF3{}, fmt.Errorf("amf3 decode array: object index outbound(expected < %d, got %d)", len(ed.deObjTable), index)
		}

		// 转换类型
		res, ok := ed.deObjTable[index].(EcmaArrayAMF3)
		if !ok {
			return EcmaArrayAMF3{}, fmt.Errorf("amf3 decode array: unable to extract array from object references")
		}

		return res, err
	}

	// 按数值取值
	var key string
	if key, err = ed.decodeString(r); err != nil {
		return EcmaArrayAMF3{}, err
	}

	// 初始化数组
	var result EcmaArrayAMF3
	ed.deObjTable = append(ed.deObjTable, result)
	result.associative = make(Object)
	result.dense = make(Array, refVal)

	// list of assoicative array, terminated by an empty string
	for key != "" {
		// 读取key对应的值
		if result.associative[key], err = ed.Decode(r); err != nil {
			return EcmaArrayAMF3{}, err
		}

		// 读取下一个key
		if key, err = ed.decodeString(r); err != nil {
			return EcmaArrayAMF3{}, err
		}
	}

	// list of dense array
	for i := uint32(0); i < refVal; i++ {
		if result.dense[i], err = ed.Decode(r); err != nil {
			return EcmaArrayAMF3{}, err
		}
	}

	return result, nil
}

func (ed *EnDecAMF3) decodeObject(r io.Reader) (interface{}, error) {
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return nil, err
	}

	// 按引用传值
	if isRef {
		index := int(refVal)
		if index >= len(ed.deObjTable) {
			return nil, fmt.Errorf("amf3 decode object: object index outbound(expected < %d, got %d)", len(ed.deObjTable), index)
		}

		return ed.deObjTable[index], nil
	}

	// each type has traits that are cached, if the peer sent a reference
	// then we'll need to look it up and use it.
	var trait TraitAMF3

	// 读取特性表
	// refs: U29O-traits-ref: representing whether a trait reference follows
	traitIsRef := (refVal & 0x01) == 0
	if traitIsRef {
		index := int(refVal >> 1)
		if index >= len(ed.traitTable) {
			return nil, fmt.Errorf("amf3 decode object: trait index outbound(expected < %d, got %d)", len(ed.traitTable), index)
		}

		trait = ed.traitTable[index]
	} else {
		// Note: no property names are included in the trait information
		trait.Externalizable = (refVal & 0x02) == 1

		// Note: another dynamic member follows until the string-type is the empty string.
		trait.Dynamic = (refVal & 0x04) == 1

		// Note: use the empty string for anonymous classes
		if trait.Type, err = ed.decodeString(r); err != nil {
			return nil, err
		}

		// Note: traits have property keys, encoded as amf3 strings
		count := int(refVal >> 3)
		trait.Properties = make([]string, count)
		for i := 0; i < count; i++ {
			if trait.Properties[i], err = ed.decodeString(r); err != nil {
				return nil, err
			}
		}

		ed.traitTable = append(ed.traitTable, trait)
	}

	var result interface{}
	ed.deObjTable = append(ed.deObjTable, result)

	// objects can be externalizable, meaning that the system has no concrete understanding of
	// their properties or how they are encodeed. in that case, we need to find and delegate behavior
	// to the right object.
	if trait.Externalizable {
		switch trait.Type {
		case "DSA":
			// AsyncMessageExt
			if result, err = ed.decodeAsyncMessage(r); err != nil {
				return nil, err
			}
		case "DSK":
			// AcknowledgeMessageExt
			if result, err = ed.decodeAcknowledgeMessage(r); err != nil {
				return nil, err
			}
		case "flex.messaging.io.ArrayCollection":
			if result, err = ed.Decode(r); err != nil {
				return nil, err
			}

			// store an extra reference to array collection container
			ed.deObjTable = append(ed.deObjTable, result)
		default:
			// 调用外部函数
			if fn, ok := ed.externalHandlers[trait.Type]; ok {
				if result, err = fn(ed, r); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("no targer function")
			}
		}

		// 外部函数调用完即可结束
		return result, nil
	}

	// 解析Properties, Dynamic
	obj := make(Object)

	// non-externalizable objects have property keys in traits, iterate through them
	// and add the read values to the object
	for _, key := range trait.Properties {
		val, err := ed.Decode(r)
		if err != nil {
			return nil, err
		}

		obj[key] = val
	}

	// if an object is dynamic, it can have extra key/value data at the ened. in this case,
	// read keys until we get an empty one.
	if trait.Dynamic {
		for {
			// key: string
			key, err := ed.decodeString(r)
			if err != nil {
				return nil, err
			}
			if key == "" {
				break
			}

			// value
			val, err := ed.Decode(r)
			if err != nil {
				return nil, err
			}

			obj[key] = val
		}
	}

	// 获取到属性值后返回
	return obj, nil
}

func (ed *EnDecAMF3) decodeXmlDoc(r io.Reader) (string, error) {
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return "", err
	}

	// 按引用取值
	if isRef {
		index := int(refVal)
		if index >= len(ed.deObjTable) {
			return "", fmt.Errorf("amf3 decode xml: object index outbound(expected < %d, got %d)", len(ed.deObjTable), index)
		}

		// 转换类型
		result, ok := ed.deObjTable[index].(string)
		if !ok {
			return "", fmt.Errorf("amf3 decode xml: cannot coerce object reference into xml string")
		}

		return result, nil
	}

	// xml字符个数
	if refVal == 0 {
		return "", nil
	}

	// 按数值取值
	buf := make([]byte, refVal)
	if _, err = r.Read(buf); err != nil {
		return "", err
	}

	// 加入缓存
	result := string(buf)
	ed.deObjTable = append(ed.deObjTable, result)

	return result, nil
}

func (ed *EnDecAMF3) decodeByteArray(r io.Reader) ([]byte, error) {
	isRef, refVal, err := ed.decodeUint29(r)
	if err != nil {
		return nil, err
	}

	// 按引用取值(deObjTable)
	if isRef {
		index := int(refVal)
		if index >= len(ed.deObjTable) {
			return nil, fmt.Errorf("amf3 decode byte array: object index outbound(expected < %d, got %d)", len(ed.deObjTable), index)
		}

		// 转换类型
		result, ok := ed.deObjTable[index].([]byte)
		if !ok {
			return nil, fmt.Errorf("amf3 decode byte array: unable to convert object ref to bytes")
		}

		return result, nil
	}

	// 字节个数
	if refVal == 0 {
		return nil, nil
	}

	// 按数值取值
	result := make([]byte, refVal)
	if _, err := r.Read(result); err != nil {
		return nil, err
	}

	// 加入缓存
	ed.deObjTable = append(ed.deObjTable, result)

	return result, nil
}

// 根据协议, 解析 U29 的语义
func (ed *EnDecAMF3) decodeUint29(r io.Reader) (bool, uint32, error) {
	u29, err := ed.readUint29(r)
	if err != nil {
		return false, 0, err
	}

	// A variable length unsigned 29-bit integer is used for the header
	// and the first bit is flag that specifies which type of string is encoded
	isRef := u29&0x01 == 0
	refVal := u29 >> 1

	return isRef, refVal, err
}

// 将压缩的整数解压为int32类型(int32->int29)
func (ed *EnDecAMF3) readUint29(r io.Reader) (uint32, error) {
	var ret uint32

	// 前3个字节
	for i := 0; i < 3; i++ {
		b, err := readByte(r)
		if err != nil {
			return 0, err
		}

		// 数据位
		ret = (ret << 7) + uint32(b&0x7F)

		// 标志位: 该位标志着下一个字节是否是整数的一部分
		if (b & 0x80) == 0 {
			return ret, nil
		}
	}

	// 第4个字节
	b, err := readByte(r)
	if err != nil {
		return 0, err
	}
	ret = (ret << 8) + uint32(b)

	return ret, nil
}

// ======================================================================

// Encode AMF3封装器入口
func (ed *EnDecAMF3) Encode(w io.Writer, val interface{}) (int, error) {
	v := reflect.ValueOf(val)

	// 复杂类型
	switch val.(type) {
	case time.Time:
		return ed.encodeDate(w, val.(time.Time))
	case Object:
		typObj := new(TypedObject)
		typObj.Object = val.(Object)

		return ed.encodeObject(w, *typObj)
	case TypedObject:
		return ed.encodeObject(w, val.(TypedObject))
	case []byte:
		return ed.encodeByteArray(w, val.([]byte))
	}

	// 简单类型
	switch v.Kind() {
	case reflect.Invalid:
		return ed.encodeNull(w)
	case reflect.Bool:
		if v.Bool() {
			return ed.encodeTrue(w)
		}

		return ed.encodeFalse(w)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		if v.Int() >= -536870911 && v.Int() < 536870911 {
			return ed.encodeInteger(w, uint32(v.Int()))
		}

		return ed.encodeDouble(w, float64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		if v.Uint() < 536870911 {
			return ed.encodeInteger(w, uint32(v.Uint()))
		}

		return ed.encodeDouble(w, float64(v.Uint()))
	case reflect.String:
		return ed.encodeString(w, v.String())
	case reflect.Int64:
		return ed.encodeDouble(w, float64(v.Int()))
	case reflect.Uint64:
		return ed.encodeDouble(w, float64(v.Uint()))
	case reflect.Float32, reflect.Float64:
		return ed.encodeDouble(w, float64(v.Float()))
	case reflect.Array, reflect.Slice:
		arr := make(Array, v.Len())
		for i := 0; i < v.Len(); i++ {
			arr[i] = v.Index(int(i)).Interface()
		}
		return ed.encodeArray(w, arr)
	}

	return 0, fmt.Errorf("encode amf3: unsupported type %s", v.Type())
}

// EncodeBatch 批量封装AMF(底层调用Encode)
func (ed *EnDecAMF3) EncodeBatch(w io.Writer, args ...interface{}) error {
	for _, v := range args {
		if _, err := ed.Encode(w, v); err != nil {
			return err
		}
	}

	return nil
}

// ======================================================================

func (ed *EnDecAMF3) encodeUndefined(w io.Writer) (int, error) {
	if err := writeByte(w, undefinedMarkerOnAMF3); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ed *EnDecAMF3) encodeNull(w io.Writer) (int, error) {
	if err := writeByte(w, nullMarkerOnAMF3); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ed *EnDecAMF3) encodeFalse(w io.Writer) (int, error) {
	if err := writeByte(w, falseMarkerOnAMF3); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ed *EnDecAMF3) encodeTrue(w io.Writer) (int, error) {
	if err := writeByte(w, trueMarkerOnAMF3); err != nil {
		return 0, err
	}

	return 1, nil
}

func (ed *EnDecAMF3) encodeInteger(w io.Writer, val uint32) (n int, err error) {
	if err = writeByte(w, integerMarkerOnAMF3); err != nil {
		return
	}
	n++

	var m = 0
	if m, err = ed.writeUint29(w, val); err != nil {
		return
	}
	n += m

	return
}

func (ed *EnDecAMF3) encodeDouble(w io.Writer, val float64) (n int, err error) {
	if err = writeByte(w, doubleMarkerOnAMF3); err != nil {
		return
	}
	n++

	if err = binary.Write(w, binary.BigEndian, val); err != nil {
		return
	}
	n += 8

	return
}

func (ed *EnDecAMF3) encodeString(w io.Writer, val string) (n int, err error) {
	if err = writeByte(w, stringMarkerOnAMF3); err != nil {
		return
	}
	n++

	var m = 0
	if m, err = ed.encodeUtf8(w, val); err != nil {
		return
	}
	n += m

	// 简化编码器: 忽略缓存(enStrTable)
	return
}

func (ed *EnDecAMF3) encodeDate(w io.Writer, val time.Time) (n int, err error) {
	if err = writeByte(w, dateMarkerOnAMF3); err != nil {
		return
	}
	n++

	if err = writeByte(w, 0x01); err != nil {
		return
	}
	n++

	u64 := float64(val.Unix()) * 1000.0
	err = binary.Write(w, binary.BigEndian, &u64)
	if err != nil {
		return
	}
	n += 8

	return
}

func (ed *EnDecAMF3) encodeArray(w io.Writer, val Array) (n int, err error) {
	if err = writeByte(w, arrayMarkerOnAMF3); err != nil {
		return
	}
	n++

	var m int
	if m, err = ed.encodeUint29(w, false, uint32(len(val))); err != nil {
		return
	}
	n += m

	m, err = ed.encodeUtf8(w, "")
	if err != nil {
		return
	}
	n += m

	for _, v := range val {
		m, err = ed.Encode(w, v)
		if err != nil {
			return
		}
		n += m
	}

	return
}

func (ed *EnDecAMF3) encodeObject(w io.Writer, val TypedObject) (n int, err error) {
	if err = writeByte(w, objectMarkerOnAMF3); err != nil {
		return
	}
	n++

	// 初始化特性表
	var trait TraitAMF3
	trait.Type = val.Type
	trait.Dynamic = false
	trait.Externalizable = false

	// 缓存属性表
	for k := range val.Object {
		trait.Properties = append(trait.Properties, k)
	}
	sort.Strings(trait.Properties)

	// 编码U29
	var u29 uint32 = 0x03
	if trait.Dynamic {
		u29 |= 0x02 << 2
	}
	if trait.Externalizable {
		u29 |= 0x01 << 2
	}
	u29 |= uint32(len(trait.Properties)) << 4

	// 写入U29
	var m = 0
	m, err = ed.writeUint29(w, u29)
	if err != nil {
		return
	}
	n += m

	// 写入类型
	m, err = ed.encodeUtf8(w, trait.Type)
	if err != nil {
		return
	}
	n += m

	// 写入key
	for _, prop := range trait.Properties {
		m, err = ed.encodeUtf8(w, prop)
		if err != nil {
			return
		}
		n += m
	}

	// 外部调用这里返回
	if trait.Externalizable {
		return
	}

	// 写入value
	for _, prop := range trait.Properties {
		m, err = ed.Encode(w, val.Object[prop])
		if err != nil {
			return
		}
		n += m
	}

	if trait.Dynamic {
		for k, v := range val.Object {
			// 找到对应的key
			var foundProp = false
			for _, prop := range trait.Properties {
				if prop == k {
					foundProp = true
					break
				}
			}

			// 如果找不到, 那么编码
			if !foundProp {
				if m, err = ed.encodeUtf8(w, k); err != nil {
					return
				}
				n += m

				if m, err = ed.Encode(w, v); err != nil {
					return
				}
				n += m
			}

			// 写入结束符
			m, err = ed.encodeUtf8(w, "")
			if err != nil {
				return
			}
			n += m
		}
	}

	return
}

func (ed *EnDecAMF3) encodeByteArray(w io.Writer, val []byte) (n int, err error) {
	if err = writeByte(w, byteArrayMarkerOnAMF3); err != nil {
		return
	}
	n++

	var m = 0
	if m, err = ed.encodeUint29(w, false, uint32(len(val))); err != nil {
		return
	}
	n += m

	// 编码字节
	m, err = w.Write(val)
	if err != nil {
		return
	}
	n += m

	return
}

func (ed *EnDecAMF3) encodeUtf8(w io.Writer, val string) (n int, err error) {
	var m = 0

	// 编码原始长度
	if m, err = ed.encodeUint29(w, false, uint32(len(val))); err != nil {
		return
	}
	n += m

	// 编码字符串
	m, err = w.Write([]byte(val))
	if err != nil {
		return
	}
	n += m

	return
}

func (ed *EnDecAMF3) encodeUint29(w io.Writer, isRef bool, val uint32) (int, error) {
	val <<= 1
	if !isRef {
		val |= 0x1
	}

	return ed.writeUint29(w, val)
}

func (ed *EnDecAMF3) writeUint29(w io.Writer, val uint32) (int, error) {
	var result []byte
	if val <= 0x0000007F {
		result = []byte{byte(val)}
	} else if val <= 0x00003FFF {
		result = []byte{byte(val>>7 | 0x80), byte(val & 0x7F)}
	} else if val <= 0x001FFFFF {
		result = []byte{byte(val>>14 | 0x80), byte(val>>7&0x7F | 0x80), byte(val & 0x7F)}
	} else if val <= 0x1FFFFFFF {
		result = []byte{byte(val>>22 | 0x80), byte(val>>15&0x7F | 0x80), byte(val>>8&0x7F | 0x80), byte(val)}
	} else {
		return 0, fmt.Errorf("amf3 encode: cannot encode u29 with value %d (out of range)", val)
	}

	return w.Write(result)
}
