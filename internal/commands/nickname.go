package commands

import (
	"errors"
	"net"

	"github.com/reonardoleis/hello/internal/manager"
)

type NicknameCommand struct {
}

func (c NicknameCommand) Execute(conn *net.Conn, manager *manager.Manager, args []string) error {
	user, ok := manager.Users[conn]
	if !ok {
		return errors.New("user not found")
	}

	user.Nickname = args[0]

	return nil
}

func (c NicknameCommand) Validate(args []string) bool {
	return len(args) == 1
}

func (c NicknameCommand) Description() string {
	return "Sets the user nickname. If the nickname is already in use, the system will reject the command."
}

func (c NicknameCommand) Name() string {
	return COMMAND_NICKNAME
}
