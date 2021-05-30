package tron

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	id  string
	gid gid

	// game related state
	direction  direction
	x, y       int
	boostTicks uint // number of ticks remaining on boost
	alive      bool

	conn Connection
}

type direction int

const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
)

func NewPlayer(conn Connection) *Player {
	uuid := uuid.NewString()
	return &Player{
		id:    uuid,
		conn:  conn,
		alive: false,
	}
}

type Connection interface {
	// pumps events from chan to ws conn
	send(e event)
	// returns chan which has incoming events from ws
	recv() <-chan event
}

func NewWebsocketConnection(ws *websocket.Conn) Connection {
	w := &wsConn{
		conn:     ws,
		recvChan: make(chan event, 128),
		sendChan: make(chan event, 128),
	}
	go w.sendLoop()
	go w.recvLoop()
	return w
}

type wsConn struct {
	conn     *websocket.Conn
	recvChan chan event
	sendChan chan event
}

func (w *wsConn) send(e event) {
	w.sendChan <- e
}

func (w *wsConn) recv() <-chan event {
	return w.recvChan
}

func (w *wsConn) recvLoop() {
	defer w.conn.Close()
	for {
		_, m, err := w.conn.ReadMessage()
		if err != nil {
			close(w.recvChan)
			return
		}
		e := strings.Trim(string(m), "\n ")
		w.recvChan <- event{
			Kind: eventType(e),
		}
	}
}

func (w *wsConn) sendLoop() {
	for e := range w.sendChan {
		if err := w.conn.WriteJSON(e); err != nil {
			log.Printf("Failed to close stream: %v", err)
		}
	}
}
