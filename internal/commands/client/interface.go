package commands

import (
	"net"

	"github.com/reonardoleis/hello/internal/manager"
)

type Command interface {
	Execute(conn *net.Conn, manager *manager.ClientManager, args []string) error
}
