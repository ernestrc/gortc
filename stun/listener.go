package stun

import "net"

type Listener struct {
}

// TODO maybe should have 2 types of listeners
// or bool flag in Listener isPacket
func Listen(net, laddr string) (Listener, error) {
	panic("todo")
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (Conn, error) {
	panic("todo")
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	panic("todo")
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	panic("todo")
}
