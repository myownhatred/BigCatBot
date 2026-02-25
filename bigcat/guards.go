package bigcat

import (
	tele "gopkg.in/telebot.v4"
)

// AccessRule — access rules bitmask
type AccessRule int

const (
	RuleAnyChat AccessRule = 1 << iota
	RuleMotherShip
	RuleAdminOnly
	RuleNotBanned
)

type CommandEntry struct {
	Rule    AccessRule
	Handler func(bh *BotHandler, c tele.Context) error
}

var commandRegistry = map[string]CommandEntry{
	Timers: {Rule: RuleAnyChat, Handler: (*BotHandler).TimersHandler},
}

func (bh *BotHandler) checkAccess(c tele.Context, rule AccessRule) bool {
	if rule&RuleMotherShip != 0 {
		if c.Chat().ID != bh.comfig.ChatsAndPeps.MotherShip {
			return false
		}
	}
	if rule&RuleNotBanned != 0 {
		if user, ok := bh.brain.Users[c.Sender().ID]; ok {
			if user.IsBanned() {
				return false
			}
		}
	}
	if rule&RuleAdminOnly != 0 {
		user, ok := bh.brain.Users[c.Sender().ID]
		if !ok || !user.IsAdmin() {
			return false
		}
	}
	return true
}
