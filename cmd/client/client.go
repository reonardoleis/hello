package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/reonardoleis/hello/internal/messages"
)

var (
	nickname = ""
)

func isDebug() bool {
	return len(os.Args) > 2 && os.Args[0] == "go"
}

func getHost() string {
	host := "localhost:8080"
	if len(os.Args) < 2 || isDebug() {
		return host
	}

	return os.Args[1]
}

func main() {
	log.Println("Connecting...")
	conn, err := net.Dial("tcp", getHost())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("Connected!")

	go handle(conn)

	for {
		reader := bufio.NewReader(os.Stdin)

		text, _ := reader.ReadString('\n')

		message := messages.Message{
			Type: messages.MessageContent,
			Len:  len(text),
			Data: text,
		}

		if strings.Contains(text, "/nickname") {
			nickname = strings.Split(text, "/nickname ")[1]
			message = messages.Message{
				Type: messages.MessageNickname,
				Len:  len(nickname),
				Data: nickname,
			}
		}

		fmt.Fprintf(conn, "%s", message.Bytes())
	}
}

func handle(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		buf = bytes.Map(func(r rune) rune {
			if r == '\x00' {
				return -1
			}

			return r
		}, buf)

		message := messages.FromBytes(buf)
		message.Show()
	}
}
