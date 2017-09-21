package main

import (
	"log"
	"net"
)

const (
	maxipdatalen = 65535
	ipheaderlen  = 20
	maxpacket    = maxipdatalen - ipheaderlen
)

var handlers map[net.Addr]*handler
var sockserv net.PacketConn

func init() {
	handlers = make(map[net.Addr]*handler)
}

func (h *handler) handle(data []byte) {
	log.Printf("received %d bytes (%s) from %v\n", len(data), data, h.addr)

	switch string(data) {
	case "ping":
		n, err := sockserv.WriteTo([]byte("pong"), h.addr)
		if err != nil {
			log.Printf("error when writing to addr %v: %+v", h.addr, err)
			return
		}

		log.Printf("sent %d bytes (pong) to %v", n, h.addr)

		// cleanup handler
		handlers[h.addr] = nil
	default:
		// ignore
	}
}

type handler struct {
	addr net.Addr
}

func main() {
	var err error
	sockserv, err = net.ListenPacket("udp", ":8012")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("router listening on %v\n", sockserv.LocalAddr())

	packet := make([]byte, maxpacket)

	for {
		n, addr, err := sockserv.ReadFrom(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		if handlers[addr] == nil {
			handlers[addr] = &handler{addr}
		}

		handlers[addr].handle(packet[:n])
	}
}
