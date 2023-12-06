package messages

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Message struct {
	Type     int    `json:"type"`
	Len      int    `json:"len"`
	Command  string `json:"command"`
	Nickname string `json:"nickname"`
	Data     string `json:"data"`
}

func New(messageType int, data string) Message {
	if len(data) > 900 {
		data = data[:900]
	}

	return Message{
		Type: messageType,
		Len:  len(data),
		Data: data,
	}
}

func NewCommand(command string, data string) Message {
	if len(data) > 900 {
		data = data[:900]
	}

	return Message{
		Type:    MessageCommand,
		Len:     len(data),
		Command: command,
		Data:    data,
	}
}

func NewSystem(data string) Message {
	if len(data) > 900 {
		data = data[:900]
	}

	return Message{
		Type:     MessageSystem,
		Len:      len(data),
		Nickname: "SYSTEM",
		Data:     data,
	}
}

func (m *Message) SetNickname(nickname string) {
	if len(nickname) > 32 {
		nickname = nickname[:32]
	}

	m.Nickname = nickname
}

func (m Message) Bytes() string {
	json, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshalling message:", err)
	}
	return string(json)
}

func FromBytes(bytes []byte) Message {
	var m Message
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		log.Println("error unmarshalling message:", err)
	}

	m.Data = strings.TrimSuffix(m.Data, "\n")
	m.Data = strings.TrimSuffix(m.Data, "\r")

	m.Data = strings.ReplaceAll(m.Data, "\n", "")
	m.Data = strings.ReplaceAll(m.Data, "\r", "")

	return m
}

func (m Message) Show() {
	fmt.Println("["+m.Nickname+"]:", m.Data)
}

func (m Message) String() string {
	return m.Nickname + ": " + m.Data
}

func (m Message) Send(conn *net.Conn) {
	fmt.Fprintf(*conn, "%s", m.Bytes())
}

func (m Message) IsCommand() bool {
	return m.Type == MessageCommand
}

func (m Message) IsResult() bool {
	return m.Type == MessageResult
}
