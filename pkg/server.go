package ntwrk

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	log  log.Logger
	port int
}

const protoErr = "Unknown protocol, expected ntwrk%s\r\n"
const actionErr = "Unknown action\r\n"
const proto = "0.1"

func NewServer(port int, logger *log.Logger) *Server {
	return &Server{
		log:  *logger,
		port: port,
	}
}

// StartServer starts a network test server on `port`.
func (s *Server) Start() {
	addr := fmt.Sprintf(":%d", s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.log.Fatal(err)
	}
	s.log.Printf("Listening on %s\n", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.log.Fatal(err)
		}
		go s.handle(conn)
	}
}

// handle starts an upload or download test on the provided TCP connection.
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	remote := formatIP(conn.RemoteAddr())
	s.log.Printf("New connection from %s", remote)

	var clientProto, action string
	fmt.Fscanf(conn, protoFmt, &clientProto, &action)
	if clientProto != proto {
		msg := fmt.Sprintf(protoErr, proto)
		conn.Write([]byte(msg))
		return
	}

	switch action {
	case "echo":
		echo(conn, 0)
		s.log.Printf("Echoed %s", remote)
	case "download":
		bytes, _ := upload(conn, 0)
		s.log.Printf("Sent %d bytes to %s", bytes, remote)
	case "upload":
		bytes, _ := download(conn, 0)
		s.log.Printf("Received %d bytes from %s", bytes, remote)
	case "whoami":
		fmt.Fprintf(conn, "%s\r\n", remote)
		s.log.Printf("Identified %s", remote)
	default:
		conn.Write([]byte(actionErr))
	}
}
