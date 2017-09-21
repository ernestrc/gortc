package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

var iface = flag.String("iface", "localhost", "interface to dial to")
var port = flag.Int("port", 8012, "port to dial to")

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *iface, *port)

	log.Printf("connecting to addr %s\n", addr)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(conn)
	_, err = fmt.Fprintf(conn, "ping")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("waiting for response.. \n")

	response := make([]byte, 4)

	for {
		status, err := reader.Read(response)
		if err != nil {
			log.Fatal(err)
		}
		switch string(response) {
		case "pong":
			log.Printf("got it\n")
			return
		default:
			log.Fatal("unknown message", status)
		}
	}
}
