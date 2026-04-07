package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func checkUserAllowed(pwd string, conn net.Conn) bool {
	_, err := conn.Write([]byte(pwd))
	if err != nil {
		fmt.Println("Error writing user to server")
	}

	resp := make([]byte, 1024)
	n, err := conn.Read(resp)

	if string(resp[:n]) == "allowed" {
		return true
	} else {
		return false
	}
}

func main() {
	addr := os.Args[1]
	var pwd = ""
	allowed := false

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Enter password:")
	for i := 0; i < 3; i++ {
		fmt.Scan(&pwd)

		if checkUserAllowed(pwd, conn) {
			allowed = true
			break
		}
	}

	if !allowed {
		fmt.Println("Incorrect password")
		conn.Close()
		os.Exit(0)
	}

	go func() {
		io.Copy(os.Stdout, conn)
		conn.Close()
	}()

	io.Copy(conn, os.Stdin)
}
