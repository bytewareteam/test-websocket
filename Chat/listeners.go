package Chat

import (
	"gitlab.com/test-ws/Chat/models"
	"gitlab.com/test-ws/User"
	"gitlab.com/test-ws/websocket"
)

func Listen(manager *websocket.Manager) {
	go SubscribeRooms(manager)
	go SubscribeMessage(manager)
}

func SubscribeRooms(m *websocket.Manager) {
	// Join
	go func(l *websocket.EventListener) {
		ch := make(chan websocket.EventListenerData)
		l.Subscribe("join", ch)
		for {
			data := <-ch
			room := data.Event.Data.(map[string]interface{})["Room"].(string)
			m.ConnectionsManager.JoinRoom(room, data.Session.WsClient)
		}
	}(m.Listeners)

	// Leave
	go func(l *websocket.EventListener) {
		ch := make(chan websocket.EventListenerData)
		l.Subscribe("leave", ch)
		for {
			data := <-ch
			room := data.Event.Data.(map[string]interface{})["Room"].(string)
			m.ConnectionsManager.LeaveRoom(room, data.Session.WsClient)
		}
	}(m.Listeners)

	// Message:Room
	go func(l *websocket.EventListener) {
		ch := make(chan websocket.EventListenerData)
		l.Subscribe("leave", ch)
		for {
			dataCh := <-ch
			data := dataCh.Event.Data.(map[string]interface{}) // Extract Data
			// Joining Room
			roomName := data["Room"].(string)
			m.ConnectionsManager.JoinRoom(roomName, dataCh.Session.WsClient)

			// Fill Message
			message := models.Message{
				From:    *dataCh.Session.WsClient.Data.(*User.User),
				Room:    roomName,
				Content: data["Content"].(string),
			}
			ev := websocket.NewEvent("message:room", message)
			// Emit
			m.ConnectionsManager.EmitRoom(roomName, ev)
		}
	}(m.Listeners)
}

func SubscribeMessage(m *websocket.Manager)  {

	// Message
	go func(l *websocket.EventListener) {
		ch := make(chan websocket.EventListenerData)
		l.Subscribe("leave", ch)
		for {
			dataCh := <-ch
			data := dataCh.Event.Data.(map[string]interface{}) // Extract Data
			// Fill Message
			to := User.NewUser(data["To"].(string))
			message := models.Message{
				From:    *dataCh.Session.WsClient.Data.(*User.User),
				To:      *to,
				Content: data["Content"].(string),
			}
			ev := websocket.NewEvent("message", message)
			// Send Message
			m.ConnectionsManager.EmitClient(to.UniqKey(), ev)
		}
	}(m.Listeners)
}