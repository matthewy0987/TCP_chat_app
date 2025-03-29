package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	host := "localhost"
	port := "6666"

	if len(os.Args) > 2 {
		host = os.Args[1]
		port = os.Args[2]
	}

	connect(host, port)
}

func connect(host, port string) {
	address := net.JoinHostPort(host, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("Connected \n")

	go inputLoop(conn)
	messageLoop(conn)
}

func inputLoop(conn net.Conn) {
	stdinScanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		
		if !stdinScanner.Scan() {
			fmt.Println("DC")
			conn.Close()
			return
		}

		fmt.Fprintf(conn, "%s\n", stdinScanner.Text())
	}
}

func messageLoop(conn net.Conn) {
	connScanner := bufio.NewScanner(conn)

	for {
		// message, err := reader.ReadString('\n')
		if !connScanner.Scan() {
		  break
		}
		clearLine()
		// fmt.Print(message)
		fmt.Printf("%s\n", connScanner.Text())
		fmt.Print("> ")
	}
}

func clearLine() {
	fmt.Print("\r\033[K")
}