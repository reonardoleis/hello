package commands

const (
	COMMAND_JOIN_ROOM   = "join_room"
	COMMAND_CREATE_ROOM = "create_room"
	COMMAND_LIST_ROOMS  = "list_rooms"
	COMMAND_NICKNAME    = "nickname"
	COMMAND_HELP        = "help"
)

var (
	COMMANDS = []string{
		COMMAND_CREATE_ROOM,
		COMMAND_JOIN_ROOM,
		COMMAND_LIST_ROOMS,
		COMMAND_NICKNAME,
		COMMAND_HELP,
	}
)
