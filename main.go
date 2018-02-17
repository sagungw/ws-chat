package main

import (
	"log"

	"github.com/sagungw/ws-chat/http"
	"github.com/sagungw/ws-chat/ws"
)

func main() {
	ws.InitGlobalChannel()

	err := http.InitHTTP()
	if err != nil {
		log.Printf("error: %v", err)
	}
}
