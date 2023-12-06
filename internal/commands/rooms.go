package commands

import (
	"errors"
	"net"

	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/room"
	"github.com/reonardoleis/hello/internal/user"
)

type JoinRoomCommand struct {
}

func (c JoinRoomCommand) Execute(conn *net.Conn, manager *manager.Manager, args []string) error {
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
	return COMMAND_JOIN_ROOM
}

type CreateRoomCommand struct{}

func (c CreateRoomCommand) Execute(conn *net.Conn, manager *manager.Manager, args []string) error {
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

	return nil
}

func (c CreateRoomCommand) Validate(args []string) bool {
	return len(args) >= 1 && len(args) <= 2
}

func (c CreateRoomCommand) Description() string {
	return "Creates a room with the given name and password. If the password is not provided, the room will be public."
}

func (c CreateRoomCommand) Name() string {
	return COMMAND_CREATE_ROOM
}
