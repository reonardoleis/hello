package commands

import (
	"errors"
	"fmt"
	"net"

	"github.com/reonardoleis/hello/internal/commands/command_names"
	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
)

type NicknameCommand struct {
}

func (c NicknameCommand) Execute(conn *net.Conn, manager *manager.ServerManager, args []string) error {
	user, ok := manager.Users[conn]
	if !ok {
		return errors.New("user not found")
	}

	message := messages.NewCommand(
		command_names.COMMAND_NICKNAME,
		fmt.Sprintf("%s %s", user.Nickname, args[0]),
	)

	room := manager.FindUserRoom(conn)
	room.Broadcast(conn, message, true)

	message = messages.NewSystem(fmt.Sprintf("%s changed nickname to %s", user.Nickname, args[0]))
	room.Broadcast(conn, message)

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
	return command_names.COMMAND_NICKNAME
}
