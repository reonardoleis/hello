package ui

import (
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
	"github.com/rivo/tview"
)

var (
	Conn *net.Conn
	app  = tview.NewApplication()

	messageInput = tview.NewInputField().
			SetLabel("Message: ").
			SetFieldWidth(40).
			SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
			SetLabelColor(tview.Styles.SecondaryTextColor).
			SetPlaceholder("")

	chat = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(false). // TODO: Fix scrollable
		SetChangedFunc(func() {
			app.Draw()
		})

	userList = tview.NewTextView().
			SetRegions(true).
			SetWrap(true).
			SetWordWrap(true).
			SetScrollable(false). // TODO: Fix scrollable
			SetChangedFunc(func() {
			app.Draw()
		})

	info = tview.NewTextView().
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(false). // TODO: Fix scrollable
		SetChangedFunc(func() {
			app.Draw()
		})

	grid = tview.NewGrid().
		SetRows(20, 1, 1).
		SetColumns(100, 5).
		SetBorders(true).
		AddItem(
			chat,
			0, 0,
			1, 1,
			5, 0,
			false).
		AddItem(
			messageInput,
			1, 0,
			1, 1,
			0, 0,
			true).
		AddItem(
			userList,
			0, 1,
			2, 2,
			0, 0,
			false).
		AddItem(
			info,
			2, 0,
			1, 3,
			0, 0,
			false)
)

func Init(manager *manager.ClientManager) {
	manager.Chat = chat
	manager.MessageInput = messageInput
	manager.Info = info
	manager.UserList = userList
	messageInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			content := messageInput.GetText()
			messageInput.SetText("")

			if strings.Contains(content, "/") {
				splitted := strings.Split(content, " ")
				commandName := splitted[0]
				args := strings.Join(splitted[1:], " ")
				message := messages.NewCommand(
					commandName[1:], args,
				)

				message.Send(Conn)
				return
			}

			AddMessage("You: "+content, true)
			message := messages.New(messages.MessageContent, content)
			message.Send(Conn)
		}
	})

	chat.SetFocusFunc(func() {
		app.SetFocus(chat)
	})

	app.SetRoot(grid, true)
	app.Run()
}

func AddMessage(message string, bold bool) {
	if bold {
		message = "[green]" + message + "[-:-:-]"
	}
	chat.SetText(chat.GetText(false) + message + "\n")
}
