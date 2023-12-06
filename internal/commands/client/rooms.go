package commands

import (
	"net"

	"github.com/reonardoleis/hello/internal/manager"
)

type JoinRoomCommand struct{}

func (c JoinRoomCommand) Execute(conn *net.Conn, manager *manager.ClientManager, args []string) error {
	nickname := args[0]
	room := args[1]
	users := args[2:]

	manager.Nickname = nickname
	manager.Room = room
	manager.Users = users

	manager.UpdateUI()

	return nil
}

type LeaveRoomCommand struct{}

func (c LeaveRoomCommand) Execute(conn *net.Conn, manager *manager.ClientManager, args []string) error {
	if manager.Nickname == args[0] {
		manager.Room = ""
		manager.Users = []string{}
		manager.Nickname = ""
		manager.Chat.SetText("")
	} else {
		for i, user := range manager.Users {
			if user == args[0] {
				manager.Users = append(manager.Users[:i], manager.Users[i+1:]...)
				break
			}
		}
	}

	manager.UpdateUI()

	return nil
}

type UserJoinedCommand struct{}

func (c UserJoinedCommand) Execute(conn *net.Conn, manager *manager.ClientManager, args []string) error {
	manager.Users = append(manager.Users, args[0])
	manager.UpdateUI()

	return nil
}
