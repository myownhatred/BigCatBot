package bigcat

import (
	"Guenhwyvar/config"

	tele "gopkg.in/telebot.v3"
)

const (
	// modes
	normal = 1
	quiet  = 2
	angry  = 3
)

type BigBrain struct {
	Comfig     config.AppConfig
	UsersFlags map[int64](UserRules)
}

type ChatRules struct {
	ChatID tele.Chat
	users  []tele.User
	mode   int
}

type UserRules struct {
	MetatronChat         int64 // chat to forward to
	MetatronFordwardFlag bool  // forwarding flag
}

func (b *BigBrain) LoadComfig() {}

func NewBigBrain() *BigBrain {
	return &BigBrain{
		Comfig:     config.AppConfig{},
		UsersFlags: make(map[int64](UserRules)),
	}
}
