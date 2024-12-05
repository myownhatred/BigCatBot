package bigcat

import (
	"Guenhwyvar/config"
	dnd "Guenhwyvar/lib/DND"
	"Guenhwyvar/lib/citizen"
	freevector "Guenhwyvar/lib/vector"

	tele "gopkg.in/telebot.v4"
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
	ChatFlags   map[int64](ChatFlags)
	Users       map[int64](citizen.Citizen)
	Party       map[int64]dnd.Char
	Game        *dnd.Game
	ChatContent map[int64](ChatContent)
	VectorChan  chan string
	VectorGame  map[int64](freevector.VectorCore)
	GenTrapMap  map[int64](GeneratorTrap)
}

type ChatRules struct {
	ChatID tele.Chat
	users  []tele.User
	mode   int
}

type ChatContent struct {
	LastPicture string
	LastLine    string
}

type GeneratorTrap struct {
	UID    int64
	ChatID int64
}

type UserRules struct {
	MetatronChat         int64 // chat to forward to
	MetatronFordwardFlag bool  // forwarding flag
}

type ChatFlags struct {
	VectorGame bool
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
		ChatFlags:   make(map[int64](ChatFlags)),
		Users:       make(map[int64](citizen.Citizen)),
		Party:       make(map[int64](dnd.Char)),
		ChatContent: make(map[int64](ChatContent)),
		Game:        dnd.NewGame(),
		VectorChan:  make(chan string),
		VectorGame:  make(map[int64](freevector.VectorCore)),
		GenTrapMap:  make(map[int64](GeneratorTrap)),
	}
}

// check if we tracking this user
func (b *BigBrain) CheckCitizen(userID int64) bool {
	_, ok := b.Users[userID]
	if ok {
		return true
	} else {
		return false
	}
}

// return true if added new record and false if it existed
func (b *BigBrain) AddCitezenCache(c citizen.Citizen) bool {
	_, ok := b.Users[c.UserID]
	if ok {
		return false
	}
	b.Users[c.UserID] = c
	return true
}
