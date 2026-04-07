package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

func getServerPassword() string {
	f, err := os.Open("server_pwd.txt")
	if err != nil {
		fmt.Println("Error opening server pwd file")
	}

	hash := make([]byte, 2048)
	n, err := f.Read(hash)

	hashstr := string(hash[:n])

	return hashstr
}

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

func checkClientAllowed(conn net.Conn) {
	pwd := make([]byte, 1024)

	for i := 0; i < 3; i++ {
		n, err := conn.Read(pwd)
		if err != nil {
			fmt.Println("Error reading pwd from client")
		}

		if string(pwd[:n]) == getServerPassword() {
			conn.Write([]byte("allowed"))
			break
		} else {
			conn.Write([]byte("not allowed"))
		}
	}
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

		checkClientAllowed(conn)

		fmt.Println("Connected with client " + conn.RemoteAddr().String())

		go handle(conn)

		fmt.Println("Handling client " + conn.RemoteAddr().String())
	}
}
