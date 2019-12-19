package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/test-ws/Chat"
	"gitlab.com/test-ws/User"
	"gitlab.com/test-ws/websocket"
)

var wsManager *websocket.Manager

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Static("/", "assets")

	wsManager = websocket.NewManager()
	e.Any("/ws", func(context echo.Context) error {
		return wsManager.NewWebSocket(context.Response(), context.Request())
	})

	User.Listen(wsManager)
	Chat.Listen(wsManager)
	_ = e.Start(":9001")
}
