package room

import "github.com/reonardoleis/hello/internal/user"

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
