package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("Listen error: %s\n", err)
		return
	}
	defer listener.Close()

	conns := make(map[string]net.Conn)

	ch := make(chan string)
	go sendToClient(ch, conns)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %s\n", err)
			return
		}
		addr := conn.RemoteAddr().String()
		s := fmt.Sprintf("login %s\n", addr)
		ch <- s

		conns[addr] = conn
		go readFromClient(ch, conn)
	}
}

func readFromClient(readch chan string, conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			//break
		}
		addr := conn.RemoteAddr().String()
		if err == io.EOF {
			t := fmt.Sprintf("Logout %s\n", addr)
			readch <- string(t)
			break
		} else if err != nil {
			fmt.Printf("Read error: %s, %s\n", addr, err)
			break
		}
		s := string(buf[:n])
		readch <- s
	}
}

func sendToClient(readch chan string, conns map[string]net.Conn) {
	fmt.Println("sendToClient()")
	for {

		select {
		case data := <-readch:
			fmt.Printf("select readch %s\n", data)
			for addr, conn := range conns {
				fmt.Printf("log: %s, %s\n", data, addr)
				n, err := conn.Write([]byte(data))
				if n == 0 {
					//break
				}
				if err == io.EOF {
					fmt.Printf("conn closed: %s, %s\n", addr, err)
					delete(conns, addr)
				} else if err != nil {
					fmt.Printf("Write error: %s, %s\n", addr, err)
					delete(conns, addr)
				}
			}
		}
	}
}
