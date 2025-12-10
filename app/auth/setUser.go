package auth

import (
	"net"
)

const BlockSize = 64

func SetUser(conn net.Conn, args []string) {
	if len(args) < 4 {
		conn.Write([]byte("-ERR wrong number of arguments for 'SETUSER' command\r\n"))
		return
	}

	username := args[2]

	password := EncryptPassword(args[3])

	UserStoreMU.Lock()
	defer UserStoreMU.Unlock()

	user, exists := UserStore[username]
	if !exists {
		user = User{
			Properties: make(map[string][]string),
		}
	} else {
		user.Properties["flags"] = []string{}
		user.Properties["passwords"] = []string{
			string(password[:]),
		}
	}

	conn.Write([]byte("+OK\r\n"))
}
