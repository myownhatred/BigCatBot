package bigcat

import (
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	AnimeMawCB = "\fanimemaw"
	GrobMawCB  = "\fgrobmaw"
	FreeCB     = "\ffreemaw"
	DNDtoBar   = "\fDNDtoBar"
	DNDtoPlaza = "\fDNDtoPlaza"
	// to remove inline buttons
	Sweep = "\fsweep"
)

func CallbackHandler(c tele.Context, serv *servitor.Servitor, brain *BigBrain) error {
	cbUniq := c.Callback().Data

	// two block
	if strings.HasPrefix(cbUniq, "\ftwo") {
		args := strings.Split(cbUniq, "\ftwo")
		id, _ := strconv.Atoi(args[1])
		event, err := serv.GetTimeWithOutTimerByID(id)
		currentTime := time.Now()
		duration := currentTime.Sub(event.Time)
		days := int(duration.Hours()) / 24
		if err != nil {
			c.Send(fmt.Sprintf("при получении названия таймера случилась беда:%s", err.Error()))
		}
		serv.ResetTimer(id)
		c.Delete()
		pik, err := memser.DaysWO(days, event.Name)
		if err != nil {
			c.Send(fmt.Sprintf("при созании картинки для сбросика таймер случчилась бида:%s", err.Error()))
		}
		m := &tele.Photo{
			File:    tele.FromDisk(pik),
			Caption: fmt.Sprintf("%s сбросил таймер\n%s", c.Callback().Sender.Username, event.Name),
		}
		return c.Send(m)
	}
	// dnd stuff
	if strings.HasPrefix(cbUniq, "\fdndMeleButtons") {
		serv.Logger.Info("callback melee buttons handler",
			slog.String("callback payload:", cbUniq))
		args := strings.Split(cbUniq, "\fdndMeleButtons")
		data := strings.Split(args[1], "_")
		id, _ := strconv.ParseInt(data[0], 10, 64)
		chatID, _ := strconv.ParseInt(data[1], 10, 64)
		c.Delete()
		serv.Logger.Info("calling target buttons func",
			slog.Int64("player ID:", id),
			slog.Int64("chat ID:", chatID))
		mes, buttons, _ := DnDTargetsButtonsPriv(c, serv, brain, c.Chat().ID)
		serv.Logger.Info("combat", "sending buttons to player", id)
		return c.Send(mes, buttons)
	}
	if strings.HasPrefix(cbUniq, "\fdndAttackTarget") {
		serv.Logger.Info("callback attack handler",
			slog.String("callback payload:", cbUniq))
		args := strings.Split(cbUniq, "\fdndAttackTarget")
		data := strings.Split(args[1], "_")
		id, _ := strconv.Atoi(data[0])
		chatID, _ := strconv.ParseInt(data[1], 10, 64)
		c.Delete()
		serv.Logger.Info("calling function to calc all sheet",
			slog.Int("target ID:", id),
			slog.Int64("chat ID:", chatID))
		return DnDAttackByCallback(c, serv, brain, id, chatID)
	}

	switch cbUniq {
	case AnimeMawCB:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		opening, err := serv.GetRandomOpening("anime")
		if err != nil {
			return c.Send(fmt.Sprintf("Неполучилось с опенингом: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("Тебе выпало послушать %s - %s", opening.Description, opening.Link))
	case GrobMawCB:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		opening, err := serv.GetRandomOpening("grob")
		if err != nil {
			return c.Send(fmt.Sprintf("Неполучилось с ГРоБом: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("Тебе выпало послушать %s - %s", opening.Description, opening.Link))
	case FreeCB:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		maw, err := serv.GetFreeMaw("open")
		if err != nil {
			return c.Send(fmt.Sprintf("Неполучилось с бесплатным мавом: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("слусай %s %s", maw.Description, maw.Link))
	case DNDtoBar:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		brain.Game.SetCurrentLocation(1)
		return c.Send(brain.Game.Lookaround())
	case DNDtoPlaza:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		brain.Game.SetCurrentLocation(0)
		return c.Send(brain.Game.Lookaround())
	case Sweep:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		return nil
	default:
		return nil
	}
}
