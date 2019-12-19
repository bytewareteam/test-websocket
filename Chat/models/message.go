package models

import (
	"gitlab.com/test-ws/User"
)

type Message struct {
	From    User.User
	To      User.User
	Room    string
	Content string
}
