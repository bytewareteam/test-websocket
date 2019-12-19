package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type (
	Manager struct {
		upgrader           *websocket.Upgrader
		ConnectionsManager *ConnectionsManager
		mutex              *sync.Mutex
		Listeners          *EventListener
	}
)

func NewManager() *Manager {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &Manager{
		upgrader:           upgrader,
		ConnectionsManager: NewConnectionsManager(),
		mutex:              &sync.Mutex{},
		Listeners:          NewEventListener(),
	}
}

func (m *Manager) NewWebSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("An error occured while upgrading the connection: %v", err)
		return err
	}

	ws := NewSession(conn, m)

	go ws.Reader()
	go ws.Writer()

	ws.On("disconnect", func(event *Event) {
		if ws.WsClient != nil {
			m.ConnectionsManager.RemoveSession(ws)
		}
	})

	return nil
}
