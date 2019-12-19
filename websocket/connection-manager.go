package websocket

type (
	ClientManager struct {
		WsClient *WsClient
		Sessions MapSessions
	}
	EntityClientMap    map[interface{}]*ClientManager
	ConnectionsManager struct {
		Clients EntityClientMap
		Rooms   MapRooms
	}
)

func NewConnectionsManager() *ConnectionsManager {
	return &ConnectionsManager{Clients: make(EntityClientMap), Rooms: make(MapRooms)}
}

/**
 * Clients
 */
func (cm *ConnectionsManager) AddClient(client *WsClient, session *Session) *ConnectionsManager {
	if _, ok := cm.Clients[client.Data.UniqKey()]; !ok {
		cm.Clients[client.Data.UniqKey()] = &ClientManager{
			WsClient: client,
			Sessions: make(MapSessions, 0),
		}
	}
	cm.AddSession(client.Data.UniqKey(), session)
	return cm
}

func (cm *ConnectionsManager) RemoveClient(client *WsClient) *ConnectionsManager {
	delete(cm.Clients, client.Data.UniqKey())
	return cm
}

func (cm *ConnectionsManager) EmitClient(entityId interface{}, e *Event) {
	if client, ok := cm.Clients[entityId]; ok {
		for _, session := range client.Sessions { // Recorrer Entidades
			session.Out <- e.Raw()
		}
	}
}

/**
 * Sessions
 */
func (cm *ConnectionsManager) AddSession(entityId interface{}, session *Session) {
	cm.Clients[entityId].Sessions[session.Id] = session
}

func (cm *ConnectionsManager) RemoveSession(session *Session) {
	if client, ok := cm.Clients[session.WsClient.Data.UniqKey()]; ok {
		delete(client.Sessions, session.Id)
		if len(cm.Clients[session.WsClient.Data.UniqKey()].Sessions) == 0 {
			for roomName := range cm.Rooms {
				cm.LeaveSessionRoom(roomName, session)
			}
			cm.RemoveClient(session.WsClient)
		}
	}
}

/**
 * Rooms
 */
func (cm *ConnectionsManager) EmitRoom(name string, ev *Event) {
	for _, clients := range cm.Rooms[name].Clients {
		for _, clientSession := range clients.Sessions {
			clientSession.Out <- ev.Raw()
		}
	}
}

func (cm *ConnectionsManager) JoinRoom(name string, client *WsClient) {
	if _, ok := cm.Rooms[name]; !ok {
		cm.Rooms[name] = NewRoom(name)
	}
	if _, ok := cm.Rooms[name].Clients[client.Data.UniqKey()]; !ok {
		cm.Rooms[name].Clients[client.Data.UniqKey()] = cm.Clients[client.Data.UniqKey()]
	}
}

func (cm *ConnectionsManager) LeaveRoom(name string, client *WsClient) {
	if room, ok := cm.Rooms[name]; ok {
		if _, exists := room.Clients[client.Data.UniqKey()]; exists {
			delete(room.Clients, client.Data.UniqKey())
		}
	}
}

func (cm *ConnectionsManager) LeaveSessionRoom(name string, session *Session) {
	if room, ok := cm.Rooms[name]; ok {
		client := session.WsClient
		if _, exists := room.Clients[client.Data.UniqKey()]; exists {
			delete(room.Clients[client.Data.UniqKey()].Sessions, session.Id)
			if len(room.Clients[client.Data.UniqKey()].Sessions) == 0 {
				delete(room.Clients, client.Data.UniqKey())
			}
		}

	}
}
