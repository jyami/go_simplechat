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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "exit" {
			break
		}
		conn.Write([]byte(t))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
