package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type ConnMSG struct {
	Conn net.Conn
	Message string
}

var ch chan ConnMSG
var chSend chan ConnMSG
var ipToNickTable map[string]string
var ipToConn map[string]net.Conn


func main() {
	ln, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}

	ch = make(chan ConnMSG)
	chSend = make(chan ConnMSG)
	ipToNickTable = make(map[string]string)
	ipToConn = make(map[string]net.Conn)

	go handleNick()
	go sendMSGS()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("connection from: %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() { // reads from client here 
		msg := scanner.Text()
		msg = strings.TrimSpace(msg)
		fmt.Printf("receiver: %s || msg: %s\n",conn.RemoteAddr(), msg)

		if !isValid(msg) {
			fmt.Fprintf(conn, "invalid message\n")
		} else {
			ch <- ConnMSG{Conn: conn, Message: msg}
		}

		// fmt.Fprintf(conn, "%s\n", scanner.Text()) //echoing back to client
	}
	fmt.Printf("client %s terminated", conn.RemoteAddr()) // handle nickname removal here
	delete(ipToNickTable, conn.RemoteAddr().String()) // deleting nick from table
	delete(ipToConn, conn.RemoteAddr().String())
}

func isValid(str string) bool {
	words := strings.Fields(str)
	if len(words) == 0 {
		return false
	}

	validCommands := map[string]bool{
		"/N": true, "/NICK": true, "/L": true, "/LIST": true, "/M": true, "/MSG": true,
	}
	if !validCommands[words[0]] {
		return false
	}

	if (words[0] == "/N" || words[0] == "/NICK") && len(words) < 2 {
		return false
	}

	if (words[0] == "/M" || words[0] == "/MSG") && len(words) < 3 {
		return false
	}

	return true
}

func handleNick() {
	for msg := range ch {
		key := msg.Conn.RemoteAddr().String()
		words := strings.Fields(msg.Message)

		if words[0] == "/NICK" || words[0] == "/N" {
			nickname := words[1] // Extract nickname
			used := false
			for _, usedNick := range ipToNickTable {
				if usedNick == nickname {
					chSend <- ConnMSG{Conn: msg.Conn, Message: "Nick is taken use /LIST"}
					used = true
					continue
				}
			}
			if used == true {
				continue
			}
			ipToNickTable[key] = nickname
			ipToConn[key] = msg.Conn
			fmt.Printf("Connection from %s set their nickname to: %s\n", key, nickname)
			chSend <- ConnMSG{Conn: msg.Conn, Message: fmt.Sprintf("Nickname set to %s", nickname)}
		} else if words[0] == "/MSG" || words[0] == "/M" {
			sender, ok := ipToNickTable[key]
			if !ok { // must set nick before sending case
				chSend <- ConnMSG{Conn: msg.Conn, Message: "set Nick first with /NICK"}
				continue
			}
			recipient := words[1] // client to receive
			var recipientsIP []string
			toSend := strings.Join(words[2:], " ") // create message to send
			
			if strings.Contains(recipient, "*"){ // broadcast case
				for recveingIP := range ipToNickTable {
					if recveingIP != key {
						recipientsIP = append(recipientsIP, recveingIP)
					}
				}
			} else { // non boradcast
				recipients := strings.Split(recipient, ",")

				// gets valid ips from nicktable
				for _, recipientNick := range recipients {
					for ip, nickname := range ipToNickTable {
						if nickname == recipientNick {
							recipientsIP = append(recipientsIP, ip)
						}
					}
				}
			}
			// sending into channel
			for _, ip := range recipientsIP {
				chSend <- ConnMSG{Conn: ipToConn[ip], Message:fmt.Sprintf("%s: %s", sender, toSend)}
			}
			
			
		} else if words[0] == "/LIST" || words[0] == "/L" {
			var allNicks []string
			for _, nickname := range ipToNickTable {
				allNicks = append(allNicks, nickname)
			}
			toSend := strings.Join(allNicks, " ")
			chSend <- ConnMSG{Conn: msg.Conn, Message: toSend}
		}

		
	}
}

func sendMSGS() {
	for msg := range chSend {
		fmt.Fprintf(msg.Conn, "%s\n", msg.Message)
	}
}