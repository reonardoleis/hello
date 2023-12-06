package manager

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type ClientManager struct {
	Room     string
	Nickname string
	Users    []string

	MessageInput *tview.InputField
	Chat         *tview.TextView
	UserList     *tview.TextView
	Info         *tview.TextView
}

func (c ClientManager) UpdateUI() {
	c.UserList.SetText(
		strings.Join(c.Users, "\n"),
	)

	if c.Nickname != "" {
		c.Info.SetText(
			fmt.Sprintf(
				"Room: %s\t\tNickname: %s",
				c.Room,
				c.Nickname,
			),
		)
	} else {
		c.Info.Clear()
	}
}
