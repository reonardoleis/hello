package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/faiface/beep"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
	commands "github.com/reonardoleis/hello/internal/commands/server"
	"github.com/reonardoleis/hello/internal/manager"
	"github.com/reonardoleis/hello/internal/messages"
	"github.com/reonardoleis/hello/internal/user"
	"github.com/reonardoleis/hello/internal/utils"
)

var (
	m                                    = &sync.Mutex{}
	serverManager *manager.ServerManager = manager.New()
)

func isDebug() bool {
	return len(os.Args) > 2 && os.Args[0] == "go"
}

func getPort() string {
	if len(os.Args) < 2 || isDebug() {
		return ":8080"
	}

	return fmt.Sprintf(":%s", os.Args[1])
}

func main() {
	initPlayer(beep.SampleRate(48000), 48000*4)
	listener, err := net.Listen("tcp", getPort())
	if err != nil {
		panic(err)
	}

	go sender()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		serverManager.Connections = append(serverManager.Connections, &conn)
		serverManager.Users[&conn] = &user.User{
			Nickname: fmt.Sprintf("Guest_%d", len(serverManager.Connections)),
			Conn:     &conn,
		}
		go handle(&conn)
	}
}

func sender() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		for _, conn := range serverManager.Connections {
			fmt.Fprintf(*conn, text+"\n")
		}
	}
}

var (
	mu      sync.Mutex
	mixer   beep.Mixer
	samples [][2]float64
	buf     []byte
	context *oto.Context
	player  *oto.Player
	done    chan struct{}
)

func Lock() {
	mu.Lock()
}

func Unlock() {
	mu.Unlock()
}

func Clear() {
	mu.Lock()
	mixer.Clear()
	mu.Unlock()
}

func update() {
	mu.Lock()
	mixer.Stream(samples)
	mu.Unlock()

	for i := range samples {
		for c := range samples[i] {
			val := samples[i][c]
			if val < -1 {
				val = -1
			}
			if val > +1 {
				val = +1
			}
			valInt16 := int16(val * (1<<15 - 1))
			low := byte(valInt16)
			high := byte(valInt16 >> 8)
			buf[i*4+c*2+0] = low
			buf[i*4+c*2+1] = high
		}
	}

	player.Write(buf)
}

func Play(s ...beep.Streamer) {
	mu.Lock()
	mixer.Add(s...)
	mu.Unlock()
}

func Close() {
	if player != nil {
		if done != nil {
			done <- struct{}{}
			done = nil
		}
		player.Close()
		context.Close()
		player = nil
	}
}

func initPlayer(sampleRate beep.SampleRate, bufferSize int) error {
	mu.Lock()
	defer mu.Unlock()

	Close()

	mixer = beep.Mixer{}

	numBytes := bufferSize * 4
	samples = make([][2]float64, bufferSize)
	buf = make([]byte, numBytes)

	var err error
	context, err = oto.NewContext(int(sampleRate), 2, 2, numBytes)
	if err != nil {
		return errors.Wrap(err, "failed to initialize speaker")
	}
	player = context.NewPlayer()

	done = make(chan struct{})

	go func() {
		for {
			select {
			default:
				update()
			case <-done:
				return
			}
		}
	}()

	return nil
}

func handle(conn *net.Conn) {
	for {
		buf := make([]byte, 453696)
		n, err := (*conn).Read(buf)
		if err != nil {
			room := serverManager.FindUserRoom(conn)
			if room != nil {
				room.RemoveUser(serverManager.Users[conn])
			}

			log.Println("error reading:", err)
			return
		}

		if n == 0 {
			continue
		}

		buf = utils.SanitizeBuffer(buf)
		message := messages.FromBytes(buf)

		if message.IsAudio() {
			_, err := message.GetAudioBuffer()
			if err != nil {
				log.Println("error getting audio buffer:", err)
			}

			continue
		}

		if message.IsCommand() {
			commandArgs := strings.Split(message.Data, " ")
			if len(commandArgs) == 1 && commandArgs[0] == "" {
				commandArgs = []string{}
			}

			command := commands.GetCommand(message.Command)
			if command == nil || !command.Validate(commandArgs) {
				message := messages.NewSystem("Invalid command")
				message.Send(conn)
				continue
			}

			err := command.Execute(conn, serverManager, commandArgs)
			if err != nil {
				message := messages.NewSystem(err.Error())
				message.Send(conn)
				continue
			}

			continue
		}

		if !serverManager.IsOnRoom(conn) {
			message := messages.NewSystem("You are not on a room")
			message.Send(conn)
			continue
		}

		nickname := serverManager.Users[conn].Nickname

		message.SetNickname(nickname)
		broadcast(conn, message.Bytes())
	}
}

func broadcast(sender *net.Conn, message string) {
	room := serverManager.FindUserRoom(sender)
	for _, user := range room.Users {
		if user.Conn == sender {
			continue
		}

		fmt.Fprintf(*user.Conn, "%s", message)
	}
}
