package hub

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Hub is a websocket hub
type Hub struct {
	connectionsMx sync.RWMutex
	connections   map[*connection]struct{}
	broadcast     chan []byte
	logMx         sync.RWMutex
	log           [][]byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	fmt.Println("Creating Hub")
	H := &Hub{
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan []byte),
		connections:   make(map[*connection]struct{}),
	}

	go func() {
		for {
			msg := <-H.broadcast
			H.connectionsMx.RLock()
			for c := range H.connections {
				select {
				case c.Send <- []byte(msg):
				case <-time.After(1 * time.Second):
					log.Printf("Shutting Down Connection")
					H.removeConnection(c)
				}
			}
			H.connectionsMx.RUnlock()
		}
	}()
	return H
}

func (H *Hub) addConnection(conn *connection) {
	fmt.Println("Adding Connection")
	H.connectionsMx.Lock()
	defer H.connectionsMx.Unlock()
	H.connections[conn] = struct{}{}
}

func (H *Hub) removeConnection(conn *connection) {
	fmt.Println("Removing Connection")
	H.connectionsMx.Lock()
	defer H.connectionsMx.Unlock()
	if _, ok := H.connections[conn]; ok {
		delete(H.connections, conn)
		close(conn.Send)
	}
}
