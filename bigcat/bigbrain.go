package bigcat

import (
	"Guenhwyvar/config"
	dnd "Guenhwyvar/lib/DND"

	tele "gopkg.in/telebot.v3"
)

const (
	// modes
	normal = 1
	quiet  = 2
	angry  = 3
)

type BigBrain struct {
	Comfig      config.AppConfig
	UsersFlags  map[int64](UserRules)
	Party       map[int64]dnd.Char
	Game        *dnd.Game
	ChatContent map[int64](ChatContent)
}

type ChatRules struct {
	ChatID tele.Chat
	users  []tele.User
	mode   int
}

type ChatContent struct {
	LastPicture tele.Photo
	LastLine    string
}

type UserRules struct {
	MetatronChat         int64 // chat to forward to
	MetatronFordwardFlag bool  // forwarding flag
}

type Pers struct {
	Name  string
	Class string
	Title string
	Race  string
}

func (b *BigBrain) LoadComfig() {}

func NewBigBrain() *BigBrain {
	return &BigBrain{
		Comfig:      config.AppConfig{},
		UsersFlags:  make(map[int64](UserRules)),
		Party:       make(map[int64](dnd.Char)),
		ChatContent: make(map[int64](ChatContent)),
		Game:        dnd.NewGame(),
	}
}
