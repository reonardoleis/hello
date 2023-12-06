package main

import (
	"log"
	"net"
	"os"

	"github.com/reonardoleis/hello/internal/messages"
	"github.com/reonardoleis/hello/internal/ui"
	"github.com/reonardoleis/hello/internal/utils"
)

var (
	nickname = ""
	conn     net.Conn
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
	var err error
	log.Println("Connecting...")
	conn, err = net.Dial("tcp", getHost())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("Connected!")

	log.Println("Initializing UI...")

	go handle(conn)
	ui.Conn = &conn
	ui.Show()
}

func handle(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		buf = utils.SanitizeBuffer(buf)
		message := messages.FromBytes(buf)

		ui.AddMessage(message.String(), false)
	}
}
