package stun

var magicCookie = [4]byte{0x21, 0x12, 0xa4, 0x42}

const stunTypeMask uint16 = 0x0110

const (
	Binding Method = 0x001
)

const (
	Request         Class = 0x000
	Indication            = 0x010
	SuccessResponse       = 0x100
	ErrorResponse         = 0x110
)

const (
	MappedAddress     AttrType = 0x0001
	Username                   = 0x0006
	MessageIntegrity           = 0x0008
	ErrorCode                  = 0x0009
	UnknownAttributes          = 0x000A
	Realm                      = 0x0014
	Nonce                      = 0x0015
	XORMappedAddress           = 0x0020
	Software                   = 0x8022
	AlternateServer            = 0x8023
	FingerPrint                = 0x8028

	// legacy
	ResponseAddress = 0x0002
	ChangeAddress   = 0x0003
	SourceAddress   = 0x0004
	ChangedAddress  = 0x0005
	Password        = 0x0007
)

var (
	ErrTryAlternateServer StunError = StunError{300, "Try Alternate Server"}
	ErrBadRequest                   = StunError{400, "Bad Request"}
	ErrUnauthorized                 = StunError{401, "Unauthorized"}
	ErrStaleCredentials             = StunError{420, "Unknown Attribute"}
	ErrStaleNonce                   = StunError{438, "Stale Nonce"}
	ErrServerError                  = StunError{500, "Server Error"}
)
