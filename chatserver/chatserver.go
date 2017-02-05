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
		//defer conn.Close()

		fmt.Printf("new connection: %s\n", conn.RemoteAddr().String())

		conns[conn.RemoteAddr().String()] = conn
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
			fmt.Printf("conn closed: %s\n", addr)
			break
		} else if err != nil {
			fmt.Printf("Read error: %s, %s\n", addr, err)
			break
		}
		s := string(buf[:n])
		fmt.Printf("received %s, %s¥n", addr, s)
		//conn.Write(buf)
		readch <- string(buf[:n])
	}
}

func sendToClient(readch chan string, conns map[string]net.Conn) {
	fmt.Println("sendToClient()")
	for {

		select {
		case data := <-readch:
			fmt.Printf("select readch %s¥n", data)
			for addr, conn := range conns {
				fmt.Printf("conn write %s, %s¥n", addr, data)
				n, err := conn.Write([]byte(data)) // もしconnがcloseされていたらどうする？現状はerror?
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
