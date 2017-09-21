package stun

import (
	"encoding/json"
	"reflect"
	"testing"
)

// the following on-the-wire STUN packets were taken from
// informational RFC https://tools.ietf.org/html/rfc5769

var rfc5769SampleRequestBytes = []byte{
	0x00, 0x01, 0x00, 0x58, //    Request type and message length
	0x21, 0x12, 0xa4, 0x42, //    Magic cookie
	0xb7, 0xe7, 0xa7, 0x01, // }
	0xbc, 0x34, 0xd6, 0x86, // }  Transaction ID
	0xfa, 0x87, 0xdf, 0xae, // }
	0x80, 0x22, 0x00, 0x10, //    SOFTWARE attribute header
	0x53, 0x54, 0x55, 0x4e, // }
	0x20, 0x74, 0x65, 0x73, // }  User-agent...
	0x74, 0x20, 0x63, 0x6c, // }  ...name
	0x69, 0x65, 0x6e, 0x74, // }
	0x00, 0x24, 0x00, 0x04, //    PRIORITY attribute header
	0x6e, 0x00, 0x01, 0xff, //    ICE priority value
	0x80, 0x29, 0x00, 0x08, //    ICE-CONTROLLED attribute header
	0x93, 0x2f, 0xf9, 0xb1, // }  Pseudo-random tie breaker...
	0x51, 0x26, 0x3b, 0x36, // }   ...for ICE control
	0x00, 0x06, 0x00, 0x09, //    USERNAME attribute header
	0x65, 0x76, 0x74, 0x6a, // }
	0x3a, 0x68, 0x36, 0x76, // }  Username (9 bytes) and padding (3 bytes)
	0x59, 0x00, 0x00, 0x00, // }
	0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
	0x9a, 0xea, 0xa7, 0x0c, // }
	0xbf, 0xd8, 0xcb, 0x56, // }
	0x78, 0x1e, 0xf2, 0xb5, // }  HMAC-SHA1 fingerprint
	0xb2, 0xd3, 0xf2, 0x49, // }
	0xc1, 0xb5, 0x71, 0xa2, // }
	0x80, 0x28, 0x00, 0x04, //    FINGERPRINT attribute header
	0xe5, 0x7a, 0x3b, 0xcf, //    CRC32 fingerprint
}

var rfc5769SampleRequest = Message{
	Class:  Request,
	Method: Binding,
	ID: []byte{
		0x21, 0x12, 0xa4, 0x42,
		0xb7, 0xe7, 0xa7, 0x01,
		0xbc, 0x34, 0xd6, 0x86,
		0xfa, 0x87, 0xdf, 0xae,
	},
	Attr: []Attribute{
		Attribute{
			Type:  Software,
			Value: []byte("STUN test client"),
		},
		Attribute{
			// PRIORITY
			Type:  AttrType(0x24),
			Value: []byte{0x6e, 0x00, 0x01, 0xff},
		},
		Attribute{
			// ICE-CONTROLLEd
			Type: AttrType(0x80<<8 | 0x29),
			Value: []byte{
				0x93, 0x2f, 0xf9, 0xb1,
				0x51, 0x26, 0x3b, 0x36,
			},
		},
		Attribute{
			Type:  Username,
			Value: []byte("evtj:h6vY"),
		},
		Attribute{
			Type: MessageIntegrity,
			Value: []byte{
				0x9a, 0xea, 0xa7, 0x0c,
				0xbf, 0xd8, 0xcb, 0x56,
				0x78, 0x1e, 0xf2, 0xb5,
				0xb2, 0xd3, 0xf2, 0x49,
				0xc1, 0xb5, 0x71, 0xa2,
			},
		},
		Attribute{
			Type:  FingerPrint,
			Value: []byte{0xe5, 0x7a, 0x3b, 0xcf},
		},
	},
}

// 2.2.  Sample IPv4 Response
var rfc5769SampleResponseBytes = []byte{
	0x01, 0x01, 0x00, 0x3c, //     Response type and message length
	0x21, 0x12, 0xa4, 0x42, //     Magic cookie
	0xb7, 0xe7, 0xa7, 0x01, // }
	0xbc, 0x34, 0xd6, 0x86, // }  Transaction ID
	0xfa, 0x87, 0xdf, 0xae, // }
	0x80, 0x22, 0x00, 0x0b, //    SOFTWARE attribute header
	0x74, 0x65, 0x73, 0x74, // }
	0x20, 0x76, 0x65, 0x63, // }  UTF-8 server name
	0x74, 0x6f, 0x72, 0x00, // }
	0x00, 0x20, 0x00, 0x08, //    XOR-MAPPED-ADDRESS attribute header
	0x00, 0x01, 0xa1, 0x47, //    Address family (IPv4) and xor'd mapped port
	0xe1, 0x12, 0xa6, 0x43, //    Xor'd mapped IPv4 address
	0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
	0x2b, 0x91, 0xf5, 0x99, // }
	0xfd, 0x9e, 0x90, 0xc3, // }
	0x8c, 0x74, 0x89, 0xf9, // }  HMAC-SHA1 fingerprint
	0x2a, 0xf9, 0xba, 0x53, // }
	0xf0, 0x6b, 0xe7, 0xd7, // }
	0x80, 0x28, 0x00, 0x04, //    FINGERPRINT attribute header
	0xc0, 0x7d, 0x4c, 0x96, //    CRC32 fingerprint
}

var rfc5769SampleResponse = Message{
	Class:  SuccessResponse,
	Method: Binding,
	ID: []byte{
		0x21, 0x12, 0xa4, 0x42,
		0xb7, 0xe7, 0xa7, 0x01,
		0xbc, 0x34, 0xd6, 0x86,
		0xfa, 0x87, 0xdf, 0xae,
	},
	Attr: []Attribute{
		Attribute{
			Type:  Software,
			Value: []byte(string("test vector")),
		},
		Attribute{
			Type: XORMappedAddress,
			Value: []byte{
				0x00, 0x01, 0xa1, 0x47,
				0xe1, 0x12, 0xa6, 0x43,
			},
		},
		Attribute{
			Type: MessageIntegrity,
			Value: []byte{
				0x2b, 0x91, 0xf5, 0x99,
				0xfd, 0x9e, 0x90, 0xc3,
				0x8c, 0x74, 0x89, 0xf9,
				0x2a, 0xf9, 0xba, 0x53,
				0xf0, 0x6b, 0xe7, 0xd7,
			},
		},
		Attribute{
			Type:  FingerPrint,
			Value: []byte{0xc0, 0x7d, 0x4c, 0x96},
		},
	},
}

// 2.3.  Sample IPv6 Response
var rfc5769SampleResponseIPv6Bytes = []byte{
	0x01, 0x01, 0x00, 0x48, //    Response type and message length
	0x21, 0x12, 0xa4, 0x42, //    Magic cookie
	0xb7, 0xe7, 0xa7, 0x01, // }
	0xbc, 0x34, 0xd6, 0x86, // }  Transaction ID
	0xfa, 0x87, 0xdf, 0xae, // }
	0x80, 0x22, 0x00, 0x0b, //    SOFTWARE attribute header
	0x74, 0x65, 0x73, 0x74, // }
	0x20, 0x76, 0x65, 0x63, // }  UTF-8 server name
	0x74, 0x6f, 0x72, 0x00, // }
	0x00, 0x20, 0x00, 0x14, //    XOR-MAPPED-ADDRESS attribute header
	0x00, 0x02, 0xa1, 0x47, //    Address family (IPv6) and xor'd mapped port.
	0x01, 0x13, 0xa9, 0xfa, // }
	0xa5, 0xd3, 0xf1, 0x79, // }  Xor'd mapped IPv6 address
	0xbc, 0x25, 0xf4, 0xb5, // }
	0xbe, 0xd2, 0xb9, 0xd9, // }
	0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
	0xa3, 0x82, 0x95, 0x4e, // }
	0x4b, 0xe6, 0x7b, 0xf1, // }
	0x17, 0x84, 0xc9, 0x7c, // }  HMAC-SHA1 fingerprint
	0x82, 0x92, 0xc2, 0x75, // }
	0xbf, 0xe3, 0xed, 0x41, // }
	0x80, 0x28, 0x00, 0x04, //    FINGERPRINT attribute header
	0xc8, 0xfb, 0x0b, 0x4c, //    CRC32 fingerprint
}

var rfc5769SampleResponseIPv6 = Message{
	Class:  SuccessResponse,
	Method: Binding,
	ID: []byte{
		0x21, 0x12, 0xa4, 0x42,
		0xb7, 0xe7, 0xa7, 0x01,
		0xbc, 0x34, 0xd6, 0x86,
		0xfa, 0x87, 0xdf, 0xae,
	},
	Attr: []Attribute{
		Attribute{
			Type:  Software,
			Value: []byte("test vector"),
		},
		Attribute{
			Type: XORMappedAddress,
			Value: []byte{
				0x00, 0x02, 0xa1, 0x47,
				0x01, 0x13, 0xa9, 0xfa,
				0xa5, 0xd3, 0xf1, 0x79,
				0xbc, 0x25, 0xf4, 0xb5,
				0xbe, 0xd2, 0xb9, 0xd9,
			},
		},
		Attribute{
			Type: MessageIntegrity,
			Value: []byte{
				0xa3, 0x82, 0x95, 0x4e,
				0x4b, 0xe6, 0x7b, 0xf1,
				0x17, 0x84, 0xc9, 0x7c,
				0x82, 0x92, 0xc2, 0x75,
				0xbf, 0xe3, 0xed, 0x41,
			},
		},
		Attribute{
			Type:  FingerPrint,
			Value: []byte{0xc8, 0xfb, 0x0b, 0x4c},
		},
	},
}

// 2.4.  Sample Request with Long-Term Authentication
var rfc5769SampleRequestLongTermAuthBytes = []byte{
	0x00, 0x01, 0x00, 0x60, //    Request type and message length
	0x21, 0x12, 0xa4, 0x42, //    Magic cookie
	0x78, 0xad, 0x34, 0x33, // }
	0xc6, 0xad, 0x72, 0xc0, // }  Transaction ID
	0x29, 0xda, 0x41, 0x2e, // }
	0x00, 0x06, 0x00, 0x12, //    USERNAME attribute header
	0xe3, 0x83, 0x9e, 0xe3, // }
	0x83, 0x88, 0xe3, 0x83, // }
	0xaa, 0xe3, 0x83, 0x83, // }  Username value (18 bytes) and padding (2 bytes)
	0xe3, 0x82, 0xaf, 0xe3, // }
	0x82, 0xb9, 0x00, 0x00, // }
	0x00, 0x15, 0x00, 0x1c, //    NONCE attribute header
	0x66, 0x2f, 0x2f, 0x34, // }
	0x39, 0x39, 0x6b, 0x39, // }
	0x35, 0x34, 0x64, 0x36, // }
	0x4f, 0x4c, 0x33, 0x34, // }  Nonce value
	0x6f, 0x4c, 0x39, 0x46, // }
	0x53, 0x54, 0x76, 0x79, // }
	0x36, 0x34, 0x73, 0x41, // }
	0x00, 0x14, 0x00, 0x0b, //    REALM attribute header
	0x65, 0x78, 0x61, 0x6d, // }
	0x70, 0x6c, 0x65, 0x2e, // }  Realm value (11 bytes) and padding (1 byte)
	0x6f, 0x72, 0x67, 0x00, // }
	0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
	0xf6, 0x70, 0x24, 0x65, // }
	0x6d, 0xd6, 0x4a, 0x3e, // }
	0x02, 0xb8, 0xe0, 0x71, // }  HMAC-SHA1 fingerprint
	0x2e, 0x85, 0xc9, 0xa2, // }
	0x8c, 0xa8, 0x96, 0x66, // }
}

var rfc5769SampleRequestLongTermAuth = Message{
	Class:  Request,
	Method: Binding,
	ID: []byte{
		0x21, 0x12, 0xa4, 0x42,
		0x78, 0xad, 0x34, 0x33,
		0xc6, 0xad, 0x72, 0xc0,
		0x29, 0xda, 0x41, 0x2e,
	},
	Attr: []Attribute{
		Attribute{
			Type:  Username,
			Value: []byte("\u30DE\u30C8\u30EA\u30C3\u30AF\u30B9"),
		},
		Attribute{
			Type: Nonce,
			Value: []byte{
				0x66, 0x2f, 0x2f, 0x34,
				0x39, 0x39, 0x6b, 0x39,
				0x35, 0x34, 0x64, 0x36,
				0x4f, 0x4c, 0x33, 0x34,
				0x6f, 0x4c, 0x39, 0x46,
				0x53, 0x54, 0x76, 0x79,
				0x36, 0x34, 0x73, 0x41,
			},
		},
		Attribute{
			Type: Realm,
			Value: []byte{
				0x65, 0x78, 0x61, 0x6d,
				0x70, 0x6c, 0x65, 0x2e,
				0x6f, 0x72, 0x67,
			},
		},
		Attribute{
			Type: MessageIntegrity,
			Value: []byte{
				0xf6, 0x70, 0x24, 0x65,
				0x6d, 0xd6, 0x4a, 0x3e,
				0x02, 0xb8, 0xe0, 0x71,
				0x2e, 0x85, 0xc9, 0xa2,
				0x8c, 0xa8, 0x96, 0x66,
			},
		},
	},
}

// classic STUN request https://tools.ietf.org/html/rfc3489
var rfc3489SampleRequestBytes = []byte{
	0x00, 0x01, 0x00, 0x58, //    Request type and message length
	0x41, 0x22, 0x39, 0x36, // }
	0xb7, 0xe7, 0xa7, 0x01, // }  Transaction ID
	0xbc, 0x34, 0xd6, 0x86, // }
	0xfa, 0x87, 0xdf, 0xae, // }
	0x80, 0x22, 0x00, 0x10, //    SOFTWARE attribute header
	0x53, 0x54, 0x55, 0x4e, // }
	0x20, 0x74, 0x65, 0x73, // }  User-agent...
	0x74, 0x20, 0x63, 0x6c, // }  ...name
	0x69, 0x65, 0x6e, 0x74, // }
	0x00, 0x24, 0x00, 0x04, //    PRIORITY attribute header
	0x6e, 0x00, 0x01, 0xff, //    ICE priority value
	0x80, 0x29, 0x00, 0x08, //    ICE-CONTROLLED attribute header
	0x93, 0x2f, 0xf9, 0xb1, // }  Pseudo-random tie breaker...
	0x51, 0x26, 0x3b, 0x36, // }   ...for ICE control
	0x00, 0x06, 0x00, 0x09, //    USERNAME attribute header
	0x65, 0x76, 0x74, 0x6a, // }
	0x3a, 0x68, 0x36, 0x76, // }  Username (9 bytes) and padding (3 bytes)
	0x59, 0x00, 0x00, 0x00, // }
	0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
	0x9a, 0xea, 0xa7, 0x0c, // }
	0xbf, 0xd8, 0xcb, 0x56, // }
	0x78, 0x1e, 0xf2, 0xb5, // }  HMAC-SHA1 fingerprint
	0xb2, 0xd3, 0xf2, 0x49, // }
	0xc1, 0xb5, 0x71, 0xa2, // }
	0x80, 0x28, 0x00, 0x04, //    FINGERPRINT attribute header
	0xe5, 0x7a, 0x3b, 0xcf, //    CRC32 fingerprint
}

var rfc3489SampleRequest = Message{
	Class:  Request,
	Method: Binding,
	ID: []byte{
		0x41, 0x22, 0x39, 0x36,
		0xb7, 0xe7, 0xa7, 0x01,
		0xbc, 0x34, 0xd6, 0x86,
		0xfa, 0x87, 0xdf, 0xae,
	},
	Attr: []Attribute{
		Attribute{
			Type:  Software,
			Value: []byte("STUN test client"),
		},
		Attribute{
			// PRIORITY
			Type:  AttrType(0x24),
			Value: []byte{0x6e, 0x00, 0x01, 0xff},
		},
		Attribute{
			// ICE-CONTROLLEd
			Type: AttrType(0x80<<8 | 0x29),
			Value: []byte{
				0x93, 0x2f, 0xf9, 0xb1,
				0x51, 0x26, 0x3b, 0x36,
			},
		},
		Attribute{
			Type:  Username,
			Value: []byte("evtj:h6vY"),
		},
		Attribute{
			Type: MessageIntegrity,
			Value: []byte{
				0x9a, 0xea, 0xa7, 0x0c,
				0xbf, 0xd8, 0xcb, 0x56,
				0x78, 0x1e, 0xf2, 0xb5,
				0xb2, 0xd3, 0xf2, 0x49,
				0xc1, 0xb5, 0x71, 0xa2,
			},
		},
		Attribute{
			Type:  FingerPrint,
			Value: []byte{0xe5, 0x7a, 0x3b, 0xcf},
		},
	},
}

// RTCP packet for testing we correctly ignore non stun packet types
var rtcpPacket = []byte{
	0x80, 0xc8, 0x00, 0x06, 0x00, 0x00, 0x00, 0x55,
	0xce, 0xa5, 0x18, 0x3a, 0x39, 0xcc, 0x7d, 0x09,
	0x23, 0xed, 0x19, 0x07, 0x00, 0x00, 0x01, 0x56,
	0x00, 0x03, 0x73, 0x50,
}

func checkType(t *testing.T, msg *Message, expectedMethod Method, expectedClass Class) {
	if msg.Method != expectedMethod {
		t.Errorf("expected method %#3x vs %#3x", expectedMethod, msg.Method)
	}

	if msg.Class != expectedClass {
		t.Errorf("expected class %#3x vs %#3x", expectedClass, msg.Class)
	}
}

func checkStunTransactionID(t *testing.T, msg *Message, expectedID string) {
	if len(expectedID) != len(msg.ID) {
		t.Fatalf("unexpected transaction ID length: %d vs %d", len(expectedID), len(msg.ID))
	}

	if len(expectedID) == 16 && !msg.IsLegacy() {
		t.Errorf("expected IsLegacy to be true")
	} else if len(expectedID) == 12 && msg.IsLegacy() {
		t.Errorf("expected IsLegacy to be false")
	}

	if expectedID != string(msg.ID) {
		t.Errorf("expected %s found %s", expectedID, string(msg.ID))
	}
}

func TestIsStun(t *testing.T) {
	// RTCP Packet
	if IsStun(rtcpPacket[:]) {
		t.Error("IsStun should return false on rtcpPacket")
	}

	// RFC 5769 test sample
	if !IsStun(rfc5769SampleRequestBytes[:]) {
		t.Error("IsStun should return true on STUN packet")
	}

	// RFC 3489 message
	if IsStun(rfc3489SampleRequestBytes[:]) {
		t.Error("IsStun should return true on a classic STUN packet")
	}

	// RFC 3489 message
	if !IsStunCompat(rfc3489SampleRequestBytes[:]) {
		t.Error("IsStun should return true on a classic STUN packet")
	}
}

type testCase struct {
	data []byte
	msg  Message
}

var marshalTestcases = []testCase{
	{rfc5769SampleRequestBytes, rfc5769SampleRequest},
	{rfc5769SampleResponseBytes, rfc5769SampleResponse},
	{rfc5769SampleResponseIPv6Bytes, rfc5769SampleResponseIPv6},
	{rfc5769SampleRequestLongTermAuthBytes, rfc5769SampleRequestLongTermAuth},
	{rfc3489SampleRequestBytes, rfc3489SampleRequest},
}

func TestUnmarshal(t *testing.T) {

	for _, tcase := range marshalTestcases {
		msg, n, err := UnmarshalCompat(tcase.data)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(tcase.msg, msg) {
			t.Errorf("expected:\n")
			e, _ := json.MarshalIndent(tcase.msg, "", "  ")
			t.Errorf(string(e))

			t.Errorf("found:\n")
			f, _ := json.MarshalIndent(msg, "", "  ")
			t.Errorf(string(f))
		}

		if n != len(tcase.data) {
			t.Errorf("expected n to be %d instead of %d", len(tcase.data), n)
		}
	}
}

func TestMarshal(t *testing.T) {
	for _, tcase := range marshalTestcases {
		data, err := Marshal(tcase.msg)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(tcase.data, data) {
			t.Errorf("expected vs found:\n")
			e, _ := json.MarshalIndent(tcase.data, "", "  ")
			t.Errorf(string(e))
			f, _ := json.MarshalIndent(data, "", "  ")
			t.Errorf(string(f))
			t.Errorf("expected vs found bytes:\n")
			t.Errorf("%#v", tcase.data)
			t.Errorf("%#v", data)
		}
	}
}

func TestNotStunError(t *testing.T) {
	_, n, err := Unmarshal(rtcpPacket)

	if n != 0 {
		t.Errorf("expected n to be 0 instead of %d", n)
	}

	if err != ErrNoStun {
		t.Fatalf("expected error to be ErrNoStun but %s found", err)
	}

	// non-backwards compatible method is used instead
	_, _, err = Unmarshal(rfc3489SampleRequestBytes)
	if err != ErrNoStun {
		t.Fatalf("expected error to be ErrNoStun but %s found", err)
	}
}

func TestIsLegacy(t *testing.T) {
	if !rfc3489SampleRequest.IsLegacy() {
		t.Fatalf("message should me marked as legacy")
	}

	if rfc5769SampleRequest.IsLegacy() {
		t.Fatalf("message should NOT me marked as legacy")
	}
}
