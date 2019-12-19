package websocket

type (
	Entity interface {
		UniqKey() interface{}
	}
	WsClient struct {
		Id   interface{}
		Data Entity
	}
)

func NewWsClient(Id interface{}, data Entity) *WsClient {
	return &WsClient{Id: Id, Data: data}
}
