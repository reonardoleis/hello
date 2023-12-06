package room

import (
	"net"

	"github.com/reonardoleis/hello/internal/messages"
	"github.com/reonardoleis/hello/internal/user"
)

type Room struct {
	Name     string
	Password string
	Users    []*user.User
}

func (r *Room) RemoveUser(user *user.User) {
	for i, u := range r.Users {
		if u == user {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
		}
	}

}

func (r Room) NeedsPassword() bool {
	return r.Password != ""
}

func (r Room) GetNicknames() []string {
	var nicknames []string

	for _, user := range r.Users {
		nicknames = append(nicknames, user.Nickname)
	}

	return nicknames
}

func (r Room) Broadcast(sender *net.Conn, message messages.Message, toSender ...bool) {
	_toSender := false
	if len(toSender) > 0 {
		_toSender = toSender[0]
	}

	for _, user := range r.Users {
		if user.Conn == sender && !_toSender {
			continue
		}

		message.Send(user.Conn)
	}
}
