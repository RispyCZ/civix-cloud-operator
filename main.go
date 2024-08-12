package main

import (
	"io"
	"log"
	"net"
)

func main() {
	// The port on which the proxy will listen for incoming connections
	listenAddr := "0.0.0.0:25565"
	// The address to which the proxy will forward the data
	targetAddr := "mc.mcjabko.cz:25565"

	// Start listening on the specified address
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening on %s and redirecting to %s", listenAddr, targetAddr)

	for {
		// Accept an incoming connection
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(clientConn, targetAddr)
	}
}

func handleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Connect to the target address
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Error connecting to target %s: %v", targetAddr, err)
		return
	}
	defer targetConn.Close()

	// Start forwarding data between client and target
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}