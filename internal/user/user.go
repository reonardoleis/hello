package user

import (
	"net"
)

type User struct {
	Conn     *net.Conn
	Nickname string
}
