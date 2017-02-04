package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("connection error: %s\n", err)
		return
	}
	defer conn.Close()

	readbuf := make([]byte, 1024)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "exit" {
			break
		}
		conn.Write([]byte(t))

		n, err := conn.Read(readbuf)
		if n == 0 {
			break
		}
		fmt.Println(n)
		if err != nil {
			fmt.Printf("Read error: %s\n", err)
		}
		fmt.Println(string(readbuf[:n]))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
