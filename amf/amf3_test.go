package amf

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// u29 测试数据
type u29TestCase struct {
	value  uint32
	expect []byte
}

var u29TestCases = []u29TestCase{
	{1, []byte{0x01}},
	{2, []byte{0x02}},
	{127, []byte{0x7F}},
	{128, []byte{0x81, 0x00}},
	{255, []byte{0x81, 0x7F}},
	{256, []byte{0x82, 0x00}},
	{0x3FFF, []byte{0xFF, 0x7F}},
	{0x4000, []byte{0x81, 0x80, 0x00}},
	{0x7FFF, []byte{0x81, 0xFF, 0x7F}},
	{0x8000, []byte{0x82, 0x80, 0x00}},
	{0x1FFFFF, []byte{0xFF, 0xFF, 0x7F}},
	{0x200000, []byte{0x80, 0xC0, 0x80, 0x00}},
	{0x3FFFFF, []byte{0x80, 0xFF, 0xFF, 0xFF}},
	{0x400000, []byte{0x81, 0x80, 0x80, 0x00}},
	{0x0FFFFFFF, []byte{0xBF, 0xFF, 0xFF, 0xFF}},
}

func TestDecodeU29(t *testing.T) {
	at := assert.New(t)

	dec := EnDecAMF3{}
	for _, tc := range u29TestCases {
		buf := bytes.NewBuffer(tc.expect)
		n, err := dec.readUint29(buf)

		at.Equal(nil, err)
		at.Equal(tc.value, n)
	}
}

func TestEnDecAMF3_Decode_Undefined(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x00})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(nil, got)
}

func TestEnDecAMF3_Decode_Null(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x01})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(nil, got)
}

func TestEnDecAMF3_Decode_alse(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x02})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(false, got)
}

func TestEnDecAMF3_Decode_True(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x03})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(true, got)
}

func TestEnDecAMF3_Decode_Integer(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x04, 0xFF, 0xFF, 0x7F})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(uint32(2097151), got)
}

func TestEnDecAMF3_Decode_Double(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{
		0x05, 0x3f, 0xf3, 0x33,
		0x33, 0x33, 0x33, 0x33,
		0x33})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal(float64(1.2), got)
}

func TestEnDecAMF3_Decode_String(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{0x06, 0x07, 'f', 'o', 'o'})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)
	at.Equal("foo", got)
}

func TestEnDecAMF3_Decode_Array(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{
		0x09, 0x13, 0x01,
		0x06, 0x03, '1',
		0x06, 0x03, '2',
		0x06, 0x03, '3',
		0x06, 0x03, '4',
		0x06, 0x03, '5',
		0x06, 0x03, '6',
		0x06, 0x03, '7',
		0x06, 0x03, '8',
		0x06, 0x03, '9',
	})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)

	array, ok := got.(EcmaArrayAMF3)
	at.Equal(true, ok)

	expect := Array{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	at.Equal(expect, array.dense)
}

func TestEnDecAMF3_Decode_Object(t *testing.T) {
	at := assert.New(t)

	dec := new(EnDecAMF3)
	buf := bytes.NewReader([]byte{
		0x0a, 0x23, 0x1f, 'o', 'r', 'g', '.', 'a',
		'm', 'f', '.', 'A', 'S', 'C', 'l', 'a',
		's', 's', 0x07, 'b', 'a', 'z', 0x07, 'f',
		'o', 'o', 0x01, 0x06, 0x07, 'b', 'a', 'r',
	})

	got, err := dec.Decode(buf)
	at.Equal(nil, err)

	object, ok := got.(Object)
	at.Equal(true, ok)

	at.Equal("bar", object["foo"])
	at.Equal(nil, object["baz"])
}

// ======================================================================

func TestEnDecAMF3_Encode_EmptyString(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	n, err := enc.Encode(buf, "")
	at.Equal(nil, err)
	at.Equal(2, n)

	expect := []byte{stringMarkerOnAMF3, 0x01}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Undefined(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	n, err := enc.encodeUndefined(buf)
	at.Equal(nil, err)
	at.Equal(1, n)

	expect := []byte{undefinedMarkerOnAMF3}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Null(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	n, err := enc.Encode(buf, nil)
	at.Equal(nil, err)
	at.Equal(1, n)

	expect := []byte{nullMarkerOnAMF3}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_False(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	n, err := enc.Encode(buf, false)
	at.Equal(nil, err)
	at.Equal(1, n)

	expect := []byte{falseMarkerOnAMF3}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_True(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	n, err := enc.Encode(buf, true)
	at.Equal(nil, err)
	at.Equal(1, n)

	expect := []byte{trueMarkerOnAMF3}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Integer(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	for _, tc := range u29TestCases {
		buf := new(bytes.Buffer)
		_, err := enc.Encode(buf, tc.value)
		at.Equal(nil, err)

		got := buf.Bytes()
		var expect bytes.Buffer
		expect.Write([]byte{integerMarkerOnAMF3})
		expect.Write(tc.expect)
		at.Equal(expect.Bytes(), got)
	}

	buf := new(bytes.Buffer)
	n, err := enc.Encode(buf, uint32(4194303))
	at.Equal(nil, err)
	at.Equal(5, n)

	expect := []byte{integerMarkerOnAMF3, 0x80, 0xFF, 0xFF, 0xFF}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Double(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	_, err := enc.Encode(buf, float64(1.2))
	at.Equal(nil, err)

	expect := []byte{doubleMarkerOnAMF3, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_String(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	_, err := enc.Encode(buf, "foo")
	at.Equal(nil, err)

	expect := []byte{stringMarkerOnAMF3, 0x07, 'f', 'o', 'o'}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Array(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	arr := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	_, err := enc.Encode(buf, arr)
	at.Equal(nil, err)

	expect := []byte{
		arrayMarkerOnAMF3, 0x13, 0x01,
		0x06, 0x03, '1',
		0x06, 0x03, '2',
		0x06, 0x03, '3',
		0x06, 0x03, '4',
		0x06, 0x03, '5',
		0x06, 0x03, '6',
		0x06, 0x03, '7',
		0x06, 0x03, '8',
		0x06, 0x03, '9',
	}
	at.Equal(expect, buf.Bytes())
}

func TestEnDecAMF3_Encode_Object(t *testing.T) {
	at := assert.New(t)

	enc := new(EnDecAMF3)
	buf := new(bytes.Buffer)

	to := TypedObject{
		Type:   "",
		Object: make(Object),
	}
	to.Type = "org.amf.ASClass"
	to.Object["foo"] = "bar"
	to.Object["baz"] = nil

	_, err := enc.Encode(buf, to)
	at.Equal(nil, err)

	expect := []byte{
		objectMarkerOnAMF3, 0x23, 0x1f, 'o',
		'r', 'g', '.', 'a',
		'm', 'f', '.', 'A',
		'S', 'C', 'l', 'a',
		's', 's', 0x07, 'b',
		'a', 'z', 0x07, 'f',
		'o', 'o', 0x01, 0x06,
		0x07, 'b', 'a', 'r',
	}
	at.Equal(expect, buf.Bytes())
}

// ======================================================================

func TestRWUint29(t *testing.T) {
	at := assert.New(t)
	buf := new(bytes.Buffer)

	ed := NewEnDecAMF3()
	n, err := ed.writeUint29(buf, 100)
	at.Equal(nil, err)
	at.Equal(1, n)

	u29, err := ed.readUint29(buf)
	at.Equal(nil, err)
	at.Equal(uint32(100), u29)
}

func TestCodeUint29(t *testing.T) {
	at := assert.New(t)
	buf := new(bytes.Buffer)

	ed := NewEnDecAMF3()
	n, err := ed.encodeUint29(buf, false, 100)
	at.Equal(nil, err)
	at.Equal(2, n)

	isRef, refVal, err := ed.decodeUint29(buf)
	at.Equal(nil, err)
	at.False(isRef)
	at.Equal(uint32(100), refVal)

	n, err = ed.encodeUint29(buf, true, 100)
	at.Equal(nil, err)
	at.Equal(2, n)

	isRef, refVal, err = ed.decodeUint29(buf)
	at.Equal(nil, err)
	at.True(isRef)
	at.Equal(uint32(100), refVal)
}

// ======================================================================

func encodeAndDecodeAMF3(val interface{}, at *assert.Assertions) (result interface{}, err error) {
	ed := NewEnDecAMF3()

	buf := new(bytes.Buffer)
	_, err = ed.Encode(buf, val)
	at.Nil(err)

	result, err = ed.Decode(buf)
	at.Nil(err)

	return
}

func compareAMF3(val interface{}, t *testing.T) {
	at := assert.New(t)

	result, err := encodeAndDecodeAMF3(val, at)
	at.Nil(err)
	at.Equal(val, result)
}

func TestEnDecAMF3_Integer(t *testing.T) {
	compareAMF3(uint32(0), t)
	compareAMF3(uint32(1245), t)
	compareAMF3(uint32(123456), t)
}

func TestEnDecAMF3_Double(t *testing.T) {
	compareAMF3(float64(3.14159), t)
	compareAMF3(float64(1234567890), t)
	compareAMF3(float64(-12345), t)
}

func TestEnDecAMF3_String(t *testing.T) {
	compareAMF3("a pup!", t)
	compareAMF3("日本語", t)
}

func TestEnDecAMF3_Bool(t *testing.T) {
	compareAMF3(true, t)
	compareAMF3(false, t)
}

func TestEnDecAMF3_Null(t *testing.T) {
	compareAMF3(nil, t)
}

func TestEnDecAMF3_Date(t *testing.T) {
	t1 := time.Unix(time.Now().Unix(), 0).UTC() // nanoseconds discarded
	t2 := time.Date(1983, 9, 4, 12, 4, 8, 0, time.UTC)

	compareAMF3(t1, t)
	compareAMF3(t2, t)
}

func TestEnDecAMF3_Array(t *testing.T) {
	at := assert.New(t)

	obj := make(Object)
	obj["key"] = "val"

	var arr Array
	arr = append(arr, "amf")
	arr = append(arr, float64(2))
	arr = append(arr, -34.95)
	arr = append(arr, true)
	arr = append(arr, false)

	res, err := encodeAndDecodeAMF3(arr, at)
	at.Nil(err)

	result, ok := res.(EcmaArrayAMF3)
	at.True(ok)
	at.Equal(arr, result.dense)
}

func TestEnDecAMF3_ByteArray(t *testing.T) {
	at := assert.New(t)

	expect := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x00}

	res, err := encodeAndDecodeAMF3(expect, at)
	at.Nil(err)

	val, ok := res.([]byte)
	at.True(ok)
	at.Equal(expect, val)
}
