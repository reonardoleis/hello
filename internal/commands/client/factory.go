package commands

import "github.com/reonardoleis/hello/internal/commands/command_names"

func GetCommand(commandName string) Command {
	switch commandName {
	case command_names.COMMAND_JOIN_ROOM:
		return JoinRoomCommand{}
	case command_names.COMMAND_NICKNAME:
		return NicknameCommand{}
	case command_names.COMMAND_LEAVE_ROOM:
		return LeaveRoomCommand{}
	case command_names.COMMAND_USER_JOINED:
		return UserJoinedCommand{}
	}

	return nil
}
