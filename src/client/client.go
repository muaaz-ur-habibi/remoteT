package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr := os.Args[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		io.Copy(os.Stdout, conn)
		conn.Close()
	}()

	io.Copy(conn, os.Stdin)
}
