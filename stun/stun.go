package stun

import (
	"fmt"
	"net"
	"time"
)

// Conn represents a STUN connection. It imlements net.Conn and net.PacketConn
// interfaces.
type Conn struct {
	conn       net.Conn
	packetConn net.PacketConn
	ctx        Context
	pending    []Message
	client     bool
}

func context() Context {
	return Context{}
}

func isSupported(network string) bool {
	switch network {
	case "udp":
	case "udp4":
	case "udp6":
	case "tcp":
	case "tcp6":
	case "tcp4":
	case "tls-tcp":
		return false // TODO
	case "tls-udp":
		return false // TODO
	default:
		return false
	}

	return true
}

// Dial connects to the address on the named network and multiplexes
// it with STUN protocol effectively granting the connection STUN capabilities.
//
// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "tls-tcp" and "tls-udp".
//
// Check net.Dial for more information about the format of address.
func Dial(network, address string) (conn net.Conn, err error) {
	if !isSupported(network) {
		err = fmt.Errorf("unsupported network: %s", network)
		return
	}
	var underlying net.Conn
	var timeout time.Duration

	timeout, err = time.ParseDuration("39.5s")
	if err != nil {
		panic(err)
	}

	d := net.Dialer{Timeout: timeout}
	ctx := context()

	underlying, err = d.DialContext(&ctx, network, address)
	if err != nil {
		return
	}

	return &Conn{client: true, conn: underlying, ctx: ctx}, nil
}

// Read reads data from the connection and processes any STUN messages read.
// The rest of data will be copied to buf. Check net.Conn.Read for more info.
func (conn *Conn) Read(buf []byte) (n int, err error) {
	panic("todo")
}

// Write writes data to the connection and also flushes any pending STUN messages.
// Check net.Conn.Write for more info.
func (conn *Conn) Write(buf []byte) (n int, err error) {
	panic("todo")
}

// ReadFrom reads a packet from the connection. If the packet is a STUN packet,
// packet will be processed instead and n = 0 will be returned. If packet
// is not a STUN packet,
// copying the payload into b. It returns the number of
// bytes copied into b and the return address that
// was on the packet.
// ReadFrom can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetReadDeadline.
func (conn *Conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	panic("todo")
}

// WriteTo writes a packet with payload b to addr.
// WriteTo can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (conn *Conn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	panic("todo")
}

func (conn *Conn) writePending() error {
	panic("todo")
}

// Close will close the underlying connection. Note that this method should be
// called only when client has no plans to use any resources such as a mapped,
// address or relayed address that were learned through STUN requests sent over
// this connection.
// If there are any pending STUN messages to deliver, method will block until
// messages have been flushed.
func (conn *Conn) Close() error {
	if len(conn.pending) == 0 {
		return conn.Close()
	}

	if err := conn.writePending(); err != nil {
		return err
	}

	return conn.Close()
}

func (conn *Conn) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

func (conn *Conn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

func (conn *Conn) SetDeadline(t time.Time) error {
	return conn.conn.SetDeadline(t)
}

func (conn *Conn) SetReadDeadline(t time.Time) error {
	return conn.conn.SetReadDeadline(t)
}

func (conn *Conn) SetWriteDeadline(t time.Time) error {
	return conn.conn.SetWriteDeadline(t)
}

// TODO STUN API
