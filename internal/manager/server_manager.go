package manager

import (
	"net"

	"github.com/reonardoleis/hello/internal/room"
	"github.com/reonardoleis/hello/internal/user"
)

type ServerManager struct {
	Connections []*net.Conn
	Users       map[*net.Conn]*user.User
	Rooms       []*room.Room
}

func New() *ServerManager {
	return &ServerManager{
		Connections: []*net.Conn{},
		Users:       map[*net.Conn]*user.User{},
		Rooms:       []*room.Room{},
	}
}

func (m ServerManager) IsOnRoom(conn *net.Conn) bool {
	for _, room := range m.Rooms {
		for _, user := range room.Users {
			if user.Conn == conn {
				return true
			}
		}
	}

	return false
}

func (m ServerManager) FindUserRoom(conn *net.Conn) *room.Room {
	for _, room := range m.Rooms {
		for _, user := range room.Users {
			if user.Conn == conn {
				return room
			}
		}
	}

	return nil
}

func (m ServerManager) RemoveUser(conn *net.Conn) {
	delete(m.Users, conn)

	for i := range m.Connections {
		if m.Connections[i] == conn {
			m.Connections = append(m.Connections[:i], m.Connections[i+1:]...)
		}
	}
}
