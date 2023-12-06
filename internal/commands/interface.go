package commands

import (
	"net"

	"github.com/reonardoleis/hello/internal/manager"
)

type Command interface {
	Execute(conn *net.Conn, manager *manager.Manager, args []string) error
	Validate(args []string) bool
	Description() string
	Name() string
}
