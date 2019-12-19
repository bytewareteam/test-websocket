package websocket

import "encoding/json"

type (
	EventHandler    func(*Event)
	MapEventHandler map[string][]EventHandler
	Event           struct {
		Name string      `json:"event"`
		Data interface{} `json:"data"`
	}
)

func NewEvent(name string, data interface{}) *Event {
	return &Event{Name: name, Data: data}
}

func NewEventFromRaw(rawData []byte) (*Event, error) {
	ev := new(Event)
	err := json.Unmarshal(rawData, ev)
	return ev, err
}

func (e *Event) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}
