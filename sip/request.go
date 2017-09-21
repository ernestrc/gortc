package sip

import "io"

// Method represents a request method
type Method string

// methods defined in rfc3261
const (
	Register Method = "REGISTER"
	Invite          = "INVITE"
	Ack             = "ACK"
	Cancel          = "CANCEL"
	Options         = "OPTIONS"
)

// Request represents a SIP request as described in rfc3261
type Request struct {
	Method
	URI
}

// Marshal encodes the given SIP request
func MarshalRequest(req *Request, writer *io.Writer) (err error) {
	// 	writer.Write()
	return
}

// Unmarshal decodes the given SIP request
// or returns an error if there was a decoding error
func UnmarshalRequest(req *Request, reader *io.Reader) (err error) {
	if req == nil || reader == nil {
		panic("EINVAL")
	}
	return
}
