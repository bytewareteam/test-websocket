package websocket

type (
	Room struct {
		Name     string
		Clients EntityClientMap
	}
	MapRooms map[string]*Room
)

func NewRoom(name string) *Room {
	return &Room{Name: name, Clients: make(EntityClientMap)}
}
