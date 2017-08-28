package stun

import (
	"net"
	"time"
)

// Conn represents a connection with RFC-5389 STUN capabilities
type Conn struct{}

func (conn *Conn) Read(b []byte) (n int, err error) {
	panic("todo")
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	panic("todo")
}

func (conn *Conn) Close() error {
	panic("todo")
}

func (conn *Conn) LocalAddr() net.Addr {
	panic("todo")
}

func (conn *Conn) RemoteAddr() net.Addr {
	panic("todo")
}

func (conn *Conn) SetDeadline(t time.Time) error {
	panic("todo")
}

func (conn *Conn) SetReadDeadline(t time.Time) error {
	panic("todo")
}

func (conn *Conn) SetWriteDeadline(t time.Time) error {
	panic("todo")
}

func Dial(network, address string) (Conn, error) {
	panic("todo")

}

func DialTimeout(network, address string, timeout time.Duration) (Conn, error) {
	panic("todo")

}
