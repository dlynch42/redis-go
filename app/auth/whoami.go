package auth

import (
	"fmt"
	"net"
)

func WhoAmI(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ACL WHOAMI' command\r\n"))
		return
	}

	AuthenticatedUsersMU.Lock()
	username, authenticated := AuthenticatedUsers[conn]
	AuthenticatedUsersMU.Unlock()

	if !authenticated {
		conn.Write([]byte("-NOAUTH Authentication required.\r\n"))
		return
	}

	response := fmt.Sprintf("$%d\r\n%s\r\n", len(username), username)
	conn.Write([]byte(response))
}
