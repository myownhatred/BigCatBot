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
	Comfig config.AppConfig
}

type ChatRules struct {
	ChatID tele.Chat
	users  []tele.User
	mode   int
}

func (b *BigBrain) LoadComfig() {}
