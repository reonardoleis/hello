package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	commands "github.com/reonardoleis/hello/internal/commands/server"
	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
	"github.com/reonardoleis/hello/internal/user"
	"github.com/reonardoleis/hello/internal/utils"
)

var (
	m                                    = &sync.Mutex{}
	serverManager *manager.ServerManager = manager.New()
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

		serverManager.Connections = append(serverManager.Connections, &conn)
		serverManager.Users[&conn] = &user.User{
			Nickname: fmt.Sprintf("Guest_%d", len(serverManager.Connections)),
			Conn:     &conn,
		}
		go handle(&conn)
	}
}

func sender() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		for _, conn := range serverManager.Connections {
			fmt.Fprintf(*conn, text+"\n")
		}
	}
}

func handle(conn *net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := (*conn).Read(buf)
		if err != nil {
			room := serverManager.FindUserRoom(conn)
			if room != nil {
				room.RemoveUser(serverManager.Users[conn])
			}

			log.Println("error reading:", err)
			return
		}

		if n == 0 {
			continue
		}

		buf = utils.SanitizeBuffer(buf)

		message := messages.FromBytes(buf)
		fmt.Printf("%+v\n", message)
		if message.IsCommand() {
			commandArgs := strings.Split(message.Data, " ")
			if len(commandArgs) == 1 && commandArgs[0] == "" {
				commandArgs = []string{}
			}

			command := commands.GetCommand(message.Command)
			if command == nil || !command.Validate(commandArgs) {
				message := messages.NewSystem("Invalid command")
				message.Send(conn)
				continue
			}

			err := command.Execute(conn, serverManager, commandArgs)
			if err != nil {
				message := messages.NewSystem(err.Error())
				message.Send(conn)
				continue
			}

			continue
		}

		if !serverManager.IsOnRoom(conn) {
			message := messages.NewSystem("You are not on a room")
			message.Send(conn)
			continue
		}

		nickname := serverManager.Users[conn].Nickname

		message.SetNickname(nickname)
		broadcast(conn, message.Bytes())
	}
}

func broadcast(sender *net.Conn, message string) {
	room := serverManager.FindUserRoom(sender)
	for _, user := range room.Users {
		if user.Conn == sender {
			continue
		}

		fmt.Fprintf(*user.Conn, "%s", message)
	}
}
