package User

import (
	"github.com/happierall/l"
	"gitlab.com/test-ws/websocket"
)

func Listen(wsm *websocket.Manager)  {
	go func(el *websocket.EventListener) {
		channel := make(chan websocket.EventListenerData)
		el.Subscribe("authenticate", channel)
		for  {
			data := <-channel
			l.Log("Data From Listener")
			l.Debug(data.Session.Id, data.Event)
		}
	}(wsm.Listeners)
}