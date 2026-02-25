package bigcat

import (
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func (bh *BotHandler) GrokSimpleAnswerHandler(c tele.Context) (err error) {
	if c.Chat().ID != bh.comfig.ChatsAndPeps.MotherShip {
		return c.Send("тут нельзя")
	}
	if c.Sender().ID != bh.comfig.ChatsAndPeps.MisterX {
		cit := bh.brain.Users[c.Sender().ID]
		cit.GrokToks--
		bh.brain.Users[c.Sender().ID] = cit
		if cit.GrokToks <= 0 {
			return c.Send(fmt.Sprintf("вам нельзя. tokens ostalos %d", cit.GrokToks))
		}
	}

	c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "грок думает...")
	report, debug, err := bh.serv.SimpleAnswer(c.Message().Payload)
	bh.serv.Logger.Info("full grok reply",
		slog.String("reply", report),
		slog.String("debug", debug),
	)
	bh.brain.DebugString = debug
	if err != nil {
		return c.Send(fmt.Sprintf("ошибка обращения к гроку: %v", err))
	} else {
		return stringPager(report, c)
	}
}
