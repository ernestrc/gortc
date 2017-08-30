package stun

import (
	"fmt"
)

// Error as described in RFC-5389 section-18.3
type Error struct {
	code int
	msg  string
}

// AttrType as described in RFC-5389 section-18.2
type AttrType uint16

// Attribute represents a STUN attribute in a STUN message.
// After the STUN header there are zero or more attributes. Each attribute
// MUST be TLV encoded, with a 16-bit type, 16-bit length, and value.
// Each STUN attribute MUST end on a 32-bit boundary.
type Attribute struct {
	Type  AttrType
	Value []byte
}

// Class indicates whether message is a request,
// a success response, an error response, or an indication
type Class uint16

// Method is a hex number in the range 0x000 - 0xFFF.
// The encoding of STUN method into a STUN message is described in RFC-5389 Section 6.
type Method uint16

// Message represents a STUN message or otherwise known as STUN packet.
type Message struct {
	Class
	Method
	// ID represents the packet's Transaction ID
	ID []byte
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

// IsStun returns whether packet in data is a STUN packet from RFC-5389
func IsStun(data []byte) bool {
	if len(data) < 20 {
		return false
	}

	return isMagicCookie(data[4:8])
}

// IsStunCompat returns whether the packet in data is a stun packet as
// described in RFC-3489. IsStun should be used first and this function
// should be used only as a fallback in scenarios where backwards compatibility
// is required.
func IsStunCompat(data []byte) bool {
	if len(data) < 20 {
		return false
	}

	return data[0]>>6 == 0
}

// IsLegacy returns true if STUN message is RFC-3489 or false if RFC-5389
func (msg *Message) IsLegacy() bool {
	return !isMagicCookie(msg.ID[:4])
}

var (
	// ErrNoStun is returned when packet is not a STUN message
	ErrNoStun = fmt.Errorf("not a STUN message")
	// ErrMalformed is returned when packet is a STUN message but there was an
	// error parsing the attributes
	ErrMalformed = fmt.Errorf("malformed message")
	// ErrIncomplete is returned when the declared
	// message length is bigger than length of packet
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

// Marshal encodes the given STUN message in binary using network-oriented format
// as described in RFC-5389 section-6
func Marshal(msg Message) (data []byte, err error) {
	// 4 bytes for type + length
	var typeLen [4]byte
	marshaledType := uint16(msg.Method) | uint16(msg.Class)
	typeLen[0] = byte(marshaledType >> 8)
	typeLen[1] = byte(marshaledType)

	data = append(data, typeLen[:]...)

	// (magic cookie +) transaction ID
	data = append(data, msg.ID[:]...)

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

func unmarshal(data []byte) (msg Message, err error) {
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
		// magic cookie + transaction ID so we maintain backwards compatibility
		ID: data[4:20],
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

// Unmarshal decodes the given packet into an RFC-5389 STUN message
// or returns an error if there was a decoding error
func Unmarshal(data []byte) (msg Message, err error) {
	if !IsStun(data) {
		err = ErrNoStun
		return
	}

	return unmarshal(data)
}

// UnmarshalCompat decodes a a given packet into a structured STUN message,
// as described in RFC-5389 maintaining bacwkards compatibility with RFC-3489.
// Unmarshal should be prefered except for maintaing backwards compatibility.
func UnmarshalCompat(data []byte) (msg Message, err error) {
	if !IsStunCompat(data) {
		err = ErrNoStun
		return
	}

	return unmarshal(data)
}
