package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var clients = make(map[*websocket.Conn]bool)

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	clients[conn] = true
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() { delete(clients, conn) }()
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for client := range clients {
			err = client.WriteMessage(mt, message)
		}
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func sdp(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	clients[conn] = true
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() { delete(clients, conn) }()
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for client := range clients {
			if client == conn {
				continue
			}
			err = client.WriteMessage(mt, message)
		}
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/chat", echo)
	http.HandleFunc("/sdp", sdp)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
