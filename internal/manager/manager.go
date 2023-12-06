package manager

import (
	"net"

	"github.com/reonardoleis/hello/internal/room"
	"github.com/reonardoleis/hello/internal/user"
)

type Manager struct {
	Connections []*net.Conn
	Users       map[*net.Conn]*user.User
	Rooms       []*room.Room
}

func New() *Manager {
	return &Manager{
		Connections: []*net.Conn{},
		Users:       map[*net.Conn]*user.User{},
		Rooms:       []*room.Room{},
	}
}

func (m Manager) IsOnRoom(conn *net.Conn) bool {
	for _, room := range m.Rooms {
		for _, user := range room.Users {
			if user.Conn == conn {
				return true
			}
		}
	}

	return false
}

func (m Manager) FindUserRoom(conn *net.Conn) *room.Room {
	for _, room := range m.Rooms {
		for _, user := range room.Users {
			if user.Conn == conn {
				return room
			}
		}
	}

	return nil
}
