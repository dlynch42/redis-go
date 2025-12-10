package cmd

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/auth"
)

func ACL(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ACL' command\r\n"))
		return
	}

	subcommand := strings.ToUpper(args[1])

	switch subcommand {
	case "WHOAMI":
		auth.WhoAmI(conn, args)
	case "GETUSER":
		auth.GetUser(conn, args)
	case "SETUSER":
		auth.SetUser(conn, args)
	default:
		conn.Write([]byte("-ERR unknown ACL subcommand\r\n"))
	}
}
