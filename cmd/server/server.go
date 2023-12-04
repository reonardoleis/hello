package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/reonardoleis/hello/cmd/messages"
)

var (
	m         = &sync.Mutex{}
	conns     []net.Conn
	nicknames map[net.Conn]string = make(map[net.Conn]string)
)

func isDebug() bool {
	return len(os.Args) > 2 && os.Args[0] == "go"
}

func getPort() string {
	if len(os.Args) < 2 || isDebug() {
		return ":8080"
	}

	return fmt.Sprintf(":%s", os.Args[1])
}

func main() {
	listener, err := net.Listen("tcp", getPort())
	if err != nil {
		panic(err)
	}

	go sender()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		conns = append(conns, conn)
		go handle(conn)
	}
}

func sender() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		for _, conn := range conns {

			fmt.Fprintf(conn, text+"\n")
		}
	}
}

func handle(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading:", err)
			return
		}

		if n == 0 {
			continue
		}

		buf = bytes.Map(func(r rune) rune {
			if r == '\x00' {
				return -1
			}

			return r
		}, buf)

		fmt.Println(string(buf))

		message := messages.FromBytes(buf)

		if message.Type == messages.MessageNickname {
			message.Data = strings.ReplaceAll(message.Data, "\n", "")
			message.Data = strings.ReplaceAll(message.Data, "\r", "")

			exists := false
			for _, nickname := range nicknames {
				if strings.EqualFold(nickname, message.Data) {
					exists = true
				}
			}

			if exists || strings.ToUpper(message.Data) == "SYSTEM" {
				fmt.Fprintf(conn, "%s", messages.NewSystem("Username already taken.").Bytes())
				continue
			}

			nicknames[conn] = message.Data
			continue
		}

		m.Lock()
		nickname := fmt.Sprintf("Guest_%d", len(conns))
		m.Unlock()
		if nicknames[conn] != "" {
			nickname = nicknames[conn]
		} else {
			nicknames[conn] = nickname
		}

		message.SetNickname(nickname)

		broadcast(conn, message.Bytes())
	}
}

func broadcast(sender net.Conn, message string) {
	for _, conn := range conns {
		if conn == sender {
			continue
		}

		fmt.Fprintf(conn, "%s", message)
	}
}
