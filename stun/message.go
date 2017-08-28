package stun

import (
	"fmt"
)

type StunError struct {
	code int
	msg  string
}

type AttrType uint16

type Attribute struct {
	Type  AttrType
	Value []byte
}

type Class uint16

type Method uint16

type Message struct {
	Class
	Method
	// Id represents the packet's Transaction ID
	Id []byte
	// Attr represents the variable number of message attributes that a STUN message can have
	Attr []Attribute
}

func isMagicCookie(data []byte) bool {
	for i, d := range data {
		if d != magicCookie[i] {
			return false
		}
	}
	return true
}

func IsStun(data []byte) bool {
	if len(data) < 20 {
		return false
	}

	// first two bits are set to 0
	return data[0]>>2 == 0
}

func (msg *Message) IsLegacy() bool {
	return !isMagicCookie(msg.Id[:4])
}

var (
	// ErrNoStun is returned when packet is not a STUN message
	ErrNoStun = fmt.Errorf("not a STUN message")
	// ErrMalformed is returned when packet is a STUN message but there was an
	// error parsing the attributes
	ErrMalformed = fmt.Errorf("malformed message")
	// ErrIncomplete is returned when the declared message length is bigger than length of packet
	ErrIncomplete = fmt.Errorf("incomplete packet")
)

func marshalAttr(data []byte, a Attribute) []byte {
	var aheader [4]byte
	var length int

	aheader[0] = byte(a.Type >> 8)
	aheader[1] = byte(a.Type)
	length = len(a.Value)
	aheader[2] = byte(length >> 8)
	aheader[3] = byte(length)

	data = append(data, aheader[:]...)
	data = append(data, a.Value[:]...)

	// align at 32-bit boundary
	for length%4 != 0 {
		length++
		data = append(data, 0x0)
	}

	return data
}

func Marshal(msg Message) (data []byte, err error) {
	// 4 bytes for type + length
	var typeLen [4]byte
	marshaledType := uint16(msg.Method) | uint16(msg.Class)
	typeLen[0] = byte(marshaledType >> 8)
	typeLen[1] = byte(marshaledType)

	data = append(data, typeLen[:]...)

	// (magic cookie +) transaction ID
	data = append(data, msg.Id[:]...)

	for _, a := range msg.Attr {
		data = marshalAttr(data, a)
	}

	// excludes header length
	lengthField := uint16(len(data) - 20)
	data[2] = byte(lengthField >> 8)
	data[3] = byte(lengthField)

	return
}

func unmarshalAttr(data []byte) (attr Attribute, length int, err error) {
	length = int(uint16(data[2])<<8 | uint16(data[3]))
	attr.Type = AttrType(uint16(data[0])<<8 | uint16(data[1]))

	if len(data) < length {
		err = ErrMalformed
		return
	}

	data = data[4:]
	widthPad := length

	// aligned at 32-bit boundary
	for widthPad%4 != 0 {
		widthPad++
	}

	attr.Value = data[:length]

	// consumed bytes
	length = widthPad

	return
}

func Unmarshal(data []byte) (msg Message, err error) {
	if !IsStun(data) {
		err = ErrNoStun
		return
	}

	// length of the message excluding header
	if len(data[20:]) < int(uint16(data[2])<<8|uint16(data[3])) {
		err = ErrIncomplete
		return
	}

	tp := uint16(data[0])<<8 | uint16(data[1])

	// unknown method/class is delegated to handler
	msg = Message{
		Method: Method(tp & ^stunTypeMask),
		Class:  Class(tp & stunTypeMask),
		// magic cookie + transaction ID so we can maintain backwards compatibility
		Id: data[4:20],
	}

	data = data[20:]

	var attr Attribute
	var length int
	for len(data) >= 4 {
		if attr, length, err = unmarshalAttr(data); err != nil {
			return
		}
		msg.Attr = append(msg.Attr, attr)
		data = data[4+length:]
	}

	return
}
