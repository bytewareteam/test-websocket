package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

type (
	MapSessions map[string]*Session
	Session     struct {
		Id       string
		Conn     *websocket.Conn
		Out      chan []byte
		In       chan []byte
		Events   MapEventHandler
		WsClient *WsClient
		Manager  *Manager
	}
)

func NewSession(conn *websocket.Conn, manager *Manager) *Session {
	return &Session{
		Id:      xid.New().String(),
		Conn:    conn,
		Out:     make(chan []byte),
		In:      make(chan []byte),
		Events:  make(MapEventHandler),
		Manager: manager,
	}
}

func (ws *Session) Reader() {
	defer func() {
		_ = ws.Conn.Close()
	}()

	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			event := NewEvent("disconnect", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WS Message Error: %v: ", err)
				event.Name = "error"
			}
			fmt.Printf("Cerrando conexión")
			ws.triggerEvent(event)
			break
		}

		event, err := NewEventFromRaw(message)

		if err != nil {
			fmt.Printf("Error parsing Message %v:", err)
		}

		ws.triggerEvent(event)
	}
}

func (ws *Session) Writer() {
	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				_ = ws.Conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
				return
			}
			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			_ = w.Close()
		}
	}
}

func (ws *Session) triggerEvent(event *Event) {
	ws.Manager.Listeners.Next(event, ws)
	if _, ok := ws.Events[event.Name]; ok {
		for _, action := range ws.Events[event.Name] {
			action(event)
		}
	}
}

func (ws *Session) On(eventName string, action EventHandler) {
	if _, ok := ws.Events[eventName]; !ok {
		ws.Events[eventName] = make([]EventHandler, 0)
	}
	ws.Events[eventName] = append(ws.Events[eventName], action)
}
