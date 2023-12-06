package commands

func GetCommand(commandName string) Command {
	switch commandName {
	case COMMAND_CREATE_ROOM:
		return CreateRoomCommand{}
	case COMMAND_JOIN_ROOM:
		return JoinRoomCommand{}
	case COMMAND_NICKNAME:
		return NicknameCommand{}
	case COMMAND_HELP:
		return HelpCommand{}
	}

	return nil
}
