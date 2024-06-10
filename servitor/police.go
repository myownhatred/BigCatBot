package servitor

import (
	"Guenhwyvar/bringer"
	"log/slog"
)

type PoliceServ struct {
	logger *slog.Logger
	pol    bringer.Police
}

func NewPoliceServ(pol bringer.Police, logger *slog.Logger) *PoliceServ {
	return &PoliceServ{
		logger: logger,
		pol:    pol,
	}
}

func (p *PoliceServ) UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (err error) {
	return p.pol.UserDefaultCheck(UserID, username, firstname, lastname, command)
}

func (p *PoliceServ) MetatronChatAdd(ChatID int64, ChatName string) (err error) {
	return p.pol.MetatronChatAdd(ChatID, ChatName)
}

func (p *PoliceServ) MetatronChatList() (IDs []int64, ChatIDs []int64, Names []string, err error) {
	return p.pol.MetatronChatList()
}
