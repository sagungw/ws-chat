package http

import (
	"log"
	"net/http"

	"github.com/sagungw/ws-chat/ws"
)

func InitHTTP() error {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws/group", globalGroup)
	http.HandleFunc("/ws/private", privateGroup)
	http.HandleFunc("/ws/user", userList)

	webapp := http.FileServer(http.Dir("./public"))
	http.Handle("/", webapp)

	log.Println("Listening HTTP requests on port 9876")
	err := http.ListenAndServe(":9876", nil)
	return err
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func globalGroup(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer conn.Close()

	ws.GroupWS.RegisterClient(conn)

	for {
		message := &ws.Message{}
		err = conn.ReadJSON(message)
		if err != nil {
			ws.GroupWS.UnregisterClient(conn)
			log.Printf("error: %v", err)
			break
		}

		ws.GroupWS.Messages <- message
	}
}

func privateGroup(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer conn.Close()

	ws.PrivateWS.RegisterClient(conn)

	for {
		message := &ws.PrivateMessage{}
		err = conn.ReadJSON(message)
		if err != nil {
			ws.PrivateWS.UnregisterClient(conn)
			log.Printf("error: %v", err)
			break
		}

		ws.PrivateWS.Messages <- message
	}
}

func userList(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer conn.Close()

	ws.UserListWS.RegisterClient(conn)

	for {
		user := &ws.User{}
		err = conn.ReadJSON(user)
		if err != nil {
			ws.PrivateWS.UnregisterClient(conn)
			log.Printf("error: %v", err)
			break
		}

		ws.UserListWS.Users <- user
	}
}
