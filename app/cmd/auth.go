package cmd

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/auth"
)

func Auth(conn net.Conn, args []string) {
	if len(args) < 3 {
		conn.Write([]byte("-ERR wrong number of arguments for 'AUTH' command\r\n"))
		return
	}

	auth.UserStoreMU.Lock()
	defer auth.UserStoreMU.Unlock()

	username := args[1]
	user, exists := auth.UserStore[username]
	if !exists {
		conn.Write([]byte("-ERR no such user\r\n"))
		return
	}

	password := auth.EncryptPassword(args[2])
	storedPasswords := user.Properties["passwords"]
	passwordMatch := false
	for _, storedPassword := range storedPasswords {
		if password == storedPassword {
			passwordMatch = true
			break
		}
	}
	if !passwordMatch {
		conn.Write([]byte("-WRONGPASS invalid username-password pair or user is disabled\r\n"))
		return
	}

	auth.AuthenticatedUsersMU.Lock()
	auth.AuthenticatedUsers[conn] = username
	auth.AuthenticatedUsersMU.Unlock()

	conn.Write([]byte("+OK\r\n"))
}
