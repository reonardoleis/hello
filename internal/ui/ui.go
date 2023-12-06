package ui

import (
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/reonardoleis/hello/internal/messages"
	"github.com/rivo/tview"
)

var (
	Conn *net.Conn
	app  = tview.NewApplication()

	messageInput = tview.NewInputField().
			SetLabel("Message: ").
			SetFieldWidth(20).
			SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
			SetFieldTextColor(tview.Styles.GraphicsColor).
			SetLabelColor(tview.Styles.SecondaryTextColor).
			SetPlaceholder("Type your message here...")

	chat = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(false). // TODO: Fix scrollable
		SetChangedFunc(func() {
			app.Draw()
		})

	grid = tview.NewGrid().
		SetRows(20, 1).
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
			true)
)

func Show() {
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
