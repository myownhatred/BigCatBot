package bigcat

import (
	"Guenhwyvar/servitor"
	"fmt"
	"strconv"
	"strings"

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
		if err != nil {
			c.Send(fmt.Sprintf("при получении названия таймера случилась беда:%s", err.Error()))
		}
		c.Delete()
		return c.Send(fmt.Sprintf("%s сбросил таймер\n%s\nбываеть", c.Callback().Sender.Username, event.Name))
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
