package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	pongWait   int = 2
	pingPeriod int = 2
	writeWait  int = 2
)

type subscription struct {
	conn *connection
	room string
}

func (s subscription) readPump() {
	c := s.conn

	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadDeadline(time.Now().Add(time.Duration(pongWait)))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(time.Duration(pongWait)))
		return nil
	})
	for {
		_, msg, err := c.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error :%v", err)
			}
			break
		}

		m := message{msg, s.room}
		h.broadcast <- m
	}

}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(time.Duration(pingPeriod))
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.TextMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(time.Duration(writeWait)))
	return c.ws.WriteMessage(mt, payload)
}

type message struct {
	data []byte
	room string
}

type hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan message
	register   chan subscription
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}

func serveWs(w *echo.Response, r *http.Request, roomID string) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomID}
	h.register <- s
	go s.writePump()
	go s.readPump()
}
