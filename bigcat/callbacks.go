package bigcat

import (
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	AnimeMawCB = "\fanimemaw"
	GrobMawCB  = "\fgrobmaw"
	FreeCB     = "\ffreemaw"
	// to remove inline buttons
	Sweep = "\fsweep"
)

func CallbackHandler(c tele.Context, serv *servitor.Servitor) error {
	cbUniq := c.Callback().Data

	// two block
	if strings.HasPrefix(cbUniq, "\ftwo") {
		args := strings.Split(cbUniq, "\ftwo")
		id, _ := strconv.Atoi(args[1])
		serv.ResetTimer(id)
		event, err := serv.GetTimeWithOutTimerByID(id)
		currentTime := time.Now()
		duration := currentTime.Sub(event.Time)
		days := int(duration.Hours()) / 24
		if err != nil {
			c.Send(fmt.Sprintf("при получении названия таймера случилась беда:%s", err.Error()))
		}
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
	case Sweep:
		_ = c.Respond(&tele.CallbackResponse{})
		c.Delete()
		return nil
	default:
		return nil
	}
}
