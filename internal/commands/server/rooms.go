package commands

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/reonardoleis/hello/internal/commands/command_names"
	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
	"github.com/reonardoleis/hello/internal/room"
	"github.com/reonardoleis/hello/internal/user"
)

type JoinRoomCommand struct {
}

func (c JoinRoomCommand) Execute(conn *net.Conn, manager *manager.ServerManager, args []string) error {
	currentRoom := manager.FindUserRoom(conn)
	if currentRoom != nil {
		currentRoom.RemoveUser(manager.Users[conn])
	}

	for _, room := range manager.Rooms {
		if room.Name == args[0] {
			if room.NeedsPassword() {
				if len(args) == 1 {
					return errors.New("room needs password")
				}

				if room.Password != args[1] {
					return errors.New("invalid password")
				}
			}

			room.Users = append(room.Users, manager.Users[conn])
			currentUsers := strings.Join(room.GetNicknames(), " ")
			message := messages.NewCommand(
				command_names.COMMAND_JOIN_ROOM,
				fmt.Sprintf("%s %s %s", manager.Users[conn].Nickname, room.Name, currentUsers),
			)

			message.Send(conn)

			message = messages.NewCommand(
				command_names.COMMAND_USER_JOINED,
				manager.Users[conn].Nickname,
			)

			room.Broadcast(conn, message)

			return nil
		}
	}

	return errors.New("room not found")
}

func (c JoinRoomCommand) Validate(args []string) bool {
	return len(args) >= 1 && len(args) <= 2
}

func (c JoinRoomCommand) Description() string {
	return "Joins a room. If the room does not exist, the system will create it."
}

func (c JoinRoomCommand) Name() string {
	return command_names.COMMAND_JOIN_ROOM
}

type CreateRoomCommand struct{}

func (c CreateRoomCommand) Execute(conn *net.Conn, manager *manager.ServerManager, args []string) error {
	roomName := args[0]
	roomPassword := ""
	if len(args) == 2 {
		roomPassword = args[1]
	}

	room := &room.Room{
		Name:     roomName,
		Password: roomPassword,
		Users: []*user.User{
			manager.Users[conn],
		},
	}

	manager.Rooms = append(manager.Rooms, room)

	nickname := manager.Users[conn].Nickname
	message := messages.NewCommand(
		command_names.COMMAND_JOIN_ROOM,
		fmt.Sprintf("%s %s %s", nickname, room.Name, nickname),
	)

	room.Broadcast(conn, message, true)

	return nil
}

func (c CreateRoomCommand) Validate(args []string) bool {
	return len(args) >= 1 && len(args) <= 2
}

func (c CreateRoomCommand) Description() string {
	return "Creates a room with the given name and password. If the password is not provided, the room will be public."
}

func (c CreateRoomCommand) Name() string {
	return command_names.COMMAND_CREATE_ROOM
}

type LeaveRoomCommand struct{}

func (c LeaveRoomCommand) Execute(conn *net.Conn, manager *manager.ServerManager, args []string) error {
	room := manager.FindUserRoom(conn)
	if room == nil {
		return errors.New("user not in a room")
	}

	message := messages.NewCommand(
		command_names.COMMAND_LEAVE_ROOM,
		manager.Users[conn].Nickname,
	)

	room.Broadcast(conn, message, true)

	room.RemoveUser(manager.Users[conn])

	message = messages.NewSystem(fmt.Sprintf("%s left the room", manager.Users[conn].Nickname))
	room.Broadcast(conn, message, true)

	return nil
}

func (c LeaveRoomCommand) Validate(args []string) bool {
	return len(args) == 0
}

func (c LeaveRoomCommand) Description() string {
	return "Leaves the current room."
}

func (c LeaveRoomCommand) Name() string {
	return command_names.COMMAND_LEAVE_ROOM
}
