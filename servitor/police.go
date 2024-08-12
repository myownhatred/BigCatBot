package servitor

import (
	"Guenhwyvar/bringer"
	"log/slog"
	"time"
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

func (p *PoliceServ) UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (ID int64, err error) {
	return p.pol.UserDefaultCheck(UserID, username, firstname, lastname, command)
}

func (p *PoliceServ) MetatronChatAdd(ChatID int64, ChatName string) (err error) {
	return p.pol.MetatronChatAdd(ChatID, ChatName)
}

func (p *PoliceServ) MetatronChatList() (IDs []int64, ChatIDs []int64, Names []string, err error) {
	return p.pol.MetatronChatList()
}

func (p *PoliceServ) Achieves(GRID int) (IDs []int, GRIDs []int, Names []string, Ranks []int, Descrs []string, err error) {
	return p.pol.Achieves(GRID)
}

func (p *PoliceServ) UserAchs(UserID int64) (IDs []int, UserIDs []int64, AchIDs []int, Dates []time.Time, Chats []string, ChatIDs []int64, err error) {
	return p.pol.UserAchs(UserID)
}

func (p *PoliceServ) UserAchAdd(UserID int64, AID int, chat string, chatID int64) (UAID int, err error) {
	return p.pol.UserAchAdd(UserID, AID, chat, chatID)
}
