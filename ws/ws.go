package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	GroupWS    *Channel
	PrivateWS  *Channel
	UserListWS *UserWSChannel
	Users      []*User
	users      = make(map[string]bool)
	Upgrader   websocket.Upgrader
)

func InitGlobalChannel() {
	GroupWS = NewChannel()
	GroupWS.Listen()

	PrivateWS = NewChannel()
	PrivateWS.Listen()

	UserListWS = NewUserWSChannel()
	UserListWS.Listen()

	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func InsertUserIfNotExist(user *User) {
	if _, ok := users[user.Username]; !ok {
		Users = append(Users, user)
		users[user.Username] = true
	}
}
