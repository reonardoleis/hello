package commands

import "github.com/reonardoleis/hello/internal/commands/command_names"

func GetCommand(commandName string) Command {
	switch commandName {
	case command_names.COMMAND_CREATE_ROOM:
		return CreateRoomCommand{}
	case command_names.COMMAND_JOIN_ROOM:
		return JoinRoomCommand{}
	case command_names.COMMAND_NICKNAME:
		return NicknameCommand{}
	case command_names.COMMAND_HELP:
		return HelpCommand{}
	case command_names.COMMAND_LEAVE_ROOM:
		return LeaveRoomCommand{}
	}

	return nil
}
