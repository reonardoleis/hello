package commands

import (
	"net"

	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
)

type HelpCommand struct{}

func (c HelpCommand) Execute(conn *net.Conn, manager *manager.Manager, args []string) error {
	commandName := args[0]
	for _, c := range COMMANDS {
		if c == commandName {
			command := GetCommand(commandName)
			message := messages.NewSystem(command.Description())
			message.Send(conn)
			return nil
		}
	}

	message := messages.NewSystem("Command not found")
	message.Send(conn)
	return nil
}

func (c HelpCommand) Validate(args []string) bool {
	return len(args) == 1
}

func (c HelpCommand) Description() string {
	return "Shows the description of given command"
}

func (c HelpCommand) Name() string {
	return COMMAND_HELP
}
