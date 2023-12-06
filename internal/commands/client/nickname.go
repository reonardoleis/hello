package commands

import (
	"net"

	"github.com/reonardoleis/hello/internal/manager"
)

type NicknameCommand struct {
}

func (c NicknameCommand) Execute(conn *net.Conn, manager *manager.ClientManager, args []string) error {
	oldNickname := args[0]
	newNickname := args[1]

	manager.Nickname = newNickname

	for i, user := range manager.Users {
		if user == oldNickname {
			manager.Users[i] = newNickname
		}
	}

	manager.UpdateUI()

	return nil
}
