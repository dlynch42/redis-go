package auth

import (
	"net"
	"sync"
)

type User struct {
	Username   string
	Properties map[string][]string
}

var UserStore = map[string]User{
	"default": {
		Properties: map[string][]string{
			"flags": {
				"nopass",
			},
			"passwords": {},
		},
	},
}
var UserStoreMU sync.Mutex

// Authenticated USers
var AuthenticatedUsers = make(map[net.Conn]string)
var AuthenticatedUsersMU sync.Mutex
