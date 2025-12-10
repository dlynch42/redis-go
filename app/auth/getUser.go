package auth

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func GetUser(conn net.Conn, args []string) {
	// Implementation for getting user information
	if len(args) < 3 {
		conn.Write([]byte("-ERR wrong number of arguments for 'GETUSER' command\r\n"))
		return
	}

	username := args[2]

	UserStoreMU.Lock()
	user, exists := UserStore[username]
	UserStoreMU.Unlock()

	if !exists {
		conn.Write([]byte("-ERR no such user\r\n"))
		return
	}

	flags := user.Properties["flags"]
	passwords := user.Properties["passwords"]

	// Construct the response
	response := "*4\r\n"                        // 4 elements
	response += "$5\r\nflags\r\n"               // 1. "flags"
	response += resp.EncodeRESPArray(flags)     // 2. flags array
	response += "$9\r\npasswords\r\n"           // 3. "passwords"
	response += resp.EncodeRESPArray(passwords) // 4. passwords array
	conn.Write([]byte(response))
}
