package amf

import (
	"fmt"
	"io"
	"math"
)

// DSA
func (ed *EnDecAMF3) decodeAsyncMessage(r io.Reader) (result Object, err error) {
	if result, err = ed.decodeAbstractMessage(r); err != nil {
		return
	}

	if err = ed.decodeExternal(r, &result, []string{"correlationId", "correlationIdBytes"}); err != nil {
		return
	}

	return
}

// DSK
func (ed *EnDecAMF3) decodeAcknowledgeMessage(r io.Reader) (result Object, err error) {
	if result, err = ed.decodeAsyncMessage(r); err != nil {
		return
	}

	if err = ed.decodeExternal(r, &result); err != nil {
		return
	}

	return
}

// Abstract external boilerplate
func (ed *EnDecAMF3) decodeAbstractMessage(r io.Reader) (result Object, err error) {
	result = make(Object)

	if err = ed.decodeExternal(r, &result,
		[]string{"body", "clientId", "destination", "headers", "messageId", "timeStamp", "timeToLive"},
		[]string{"clientIdBytes", "messageIdBytes"}); err != nil {
		return result, fmt.Errorf("unable to decode abstract external: %s", err)
	}

	return
}

func (ed *EnDecAMF3) decodeExternal(r io.Reader, obj *Object, fieldSets ...[]string) (err error) {
	var flagSet []uint8
	if flagSet, err = ed.readFlags(r); err != nil {
		return err
	}

	var fieldNames []string
	var reservedPosition uint8
	for i, flags := range flagSet {
		if i < len(fieldSets) {
			fieldNames = fieldSets[i]
		} else {
			fieldNames = []string{}
		}

		reservedPosition = uint8(len(fieldNames))

		// 只有在flagSet里设置为1的值才会被解开
		for p, field := range fieldNames {
			flagBit := uint8(math.Exp2(float64(p)))
			if (flags & flagBit) == 1 {
				tmp, err := ed.Decode(r)
				if err != nil {
					return err
				}
				(*obj)[field] = tmp
			}
		}

		// 如果flagSet中被设置为1的位, 不存在于filedNames中, 则依旧会被解开, 并赋予默认格式的key
		if (flags >> reservedPosition) != 0 {
			for j := reservedPosition; j <= 6; j++ {
				if ((flags >> j) & 0x01) == 1 {
					tmp, err := ed.Decode(r)
					if err != nil {
						return err
					}

					field := fmt.Sprintf("extra_%d_%d", i, j)
					(*obj)[field] = tmp
				}
			}
		}
	}

	return
}

// byte的最后一位是标志位
func (ed *EnDecAMF3) readFlags(r io.Reader) (result []uint8, err error) {
	var flag byte
	for {
		if flag, err = readByte(r); err != nil {
			return
		}

		result = append(result, flag)
		if (flag & 0x80) == 0 {
			break
		}
	}

	return
}
