package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("connection error: %s\n", err)
		return
	}

	go readFromServer(conn)
	writeToServer(conn)
}

func writeToServer(conn net.Conn) {
	fmt.Println("writeToServer start")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "exit" {
			conn.Close()
			break
		}
		conn.Write([]byte(t))
	}
	fmt.Println("writeToServer end")
}

func readFromServer(conn net.Conn) {
	fmt.Println("readFromServer start")
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			//break
		}
		if err == io.EOF {
			fmt.Printf("conn closed:\n")
			break
		} else if err != nil {
			fmt.Printf("Read error: %s\n", err)
		}
		s := string(buf[:n])
		fmt.Printf("%s\n", s)
	}
	fmt.Println("readFromServer end")
}
