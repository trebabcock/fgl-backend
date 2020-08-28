package hub

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {
	Send chan []byte
	H    *Hub
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	fmt.Println("reader")
	defer wg.Done()
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		c.H.broadcast <- message
	}
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	fmt.Println("writer")
	defer wg.Done()
	for message := range c.Send {
		err := wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

// WsHandler holds the Hub
type WsHandler struct {
	H *Hub
}

func (wsh WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Open")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	c := &connection{Send: make(chan []byte, 256), H: wsh.H}
	c.H.addConnection(c)
	defer c.H.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn)
	wg.Wait()
	wsConn.Close()
	fmt.Println("Websocket closed")
}
