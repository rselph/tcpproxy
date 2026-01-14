package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {
	var (
		listenAddr  string
		forwardAddr string
		protocol    string
	)
	flag.StringVar(&listenAddr, "listen", "", "Address to listen on")
	flag.StringVar(&forwardAddr, "forward", "", "Address to forward requests to")
	flag.StringVar(&protocol, "protocol", "tcp4", "Protocol to use (tcp, tcp4, tcp6)")
	flag.Parse()

	if forwardAddr == "" {
		fmt.Println("Forward address is required")
		return
	}

	listener, err := net.Listen(protocol, listenAddr)
	if err != nil {
		fmt.Println("Listen:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on:", listener.Addr())

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept:", err)
			continue
		}
		go handleConnection(clientConn, forwardAddr)
	}
}

func handleConnection(clientConn net.Conn, forwardAddr string) {
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", forwardAddr)
	if err != nil {
		fmt.Println("Connect:", err)
		return
	}
	defer serverConn.Close()

	fmt.Println("Proxying", clientConn.RemoteAddr(), "<->", serverConn.RemoteAddr())

	var wg sync.WaitGroup
	wg.Add(2)
	go halfPipe(&wg, clientConn, serverConn)
	go halfPipe(&wg, serverConn, clientConn)
	wg.Wait()

	fmt.Println("Closed", clientConn.RemoteAddr(), "<->", serverConn.RemoteAddr())
}

func halfPipe(wg *sync.WaitGroup, src net.Conn, dst net.Conn) {
	defer wg.Done()
	buf := make([]byte, 65536)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
				fmt.Println("Read:", err)
			}
			dst.Close()
			return
		}
		_, err = dst.Write(buf[:n])
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
				fmt.Println("Write:", err)
			}
			src.Close()
			return
		}
	}
}
