package command_names

const (
	COMMAND_JOIN_ROOM   = "join_room"
	COMMAND_CREATE_ROOM = "create_room"
	COMMAND_LIST_ROOMS  = "list_rooms"
	COMMAND_NICKNAME    = "nickname"
	COMMAND_HELP        = "help"
	COMMAND_LEAVE_ROOM  = "leave_room"
	COMMAND_USER_JOINED = "user_joined"
)

var (
	COMMANDS = []string{
		COMMAND_CREATE_ROOM,
		COMMAND_JOIN_ROOM,
		COMMAND_LIST_ROOMS,
		COMMAND_NICKNAME,
		COMMAND_HELP,
		COMMAND_LEAVE_ROOM,
		COMMAND_USER_JOINED,
	}
)
